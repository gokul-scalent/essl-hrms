import React, { useEffect, useState } from "react";
import {
  FormControlLabel,
  Grid,
  Icon,
  ListItemIcon,
  ListItemText,
  MenuItem,
  Radio,
  RadioGroup,
  Select,
} from "@mui/material";
import {
  ArrowDropDown,
  Block,
  Business,
  Email,
  LocationCity,
  Person,
  Public,
  Work,
} from "@mui/icons-material";
import colors from "assets/theme/base/colors";
import { useNotify } from "components/CommonComponent/NotificationProvider";
import MDBox from "components/MDBox";
import MDButton from "components/MDButton";
import MDInput from "components/MDInput";
import MDTypography from "components/MDTypography";
import RequiredLabel from "components/CommonComponent/RequiredLabel";
import CenterPopup from "components/CommonComponent/CenterPopup";
import {
  radioLabelSx,
  leadCsvStyles,
  leadListUiStyles,
} from "components/CommonComponent/CommonFunction";
import { selectSx } from "components/CommonComponent/CommonFunction";
import { REGEX } from "components/common/constant";
import { addLeadList, getLeadsList } from "actions/leads";
import { useDispatch } from "react-redux";

const initialFields = {
  type: { value: "csv", error: "" },
  email: { value: "", error: "" },
  firstName: { value: "", error: "" },
  lastName: { value: "", error: "" },
  jobTitle: { value: "", error: "" },
  company: { value: "", error: "" },
  // country: { value: "", error: "" },
  // city: { value: "", error: "" },
  industry: { value: "", error: "" },
};

const CSV_MAPPING_FIELDS = [
  { value: "IGNORE", label: "Ignore", icon: Block },
  { value: "EMAIL", label: "Email", icon: Email },
  { value: "FIRST_NAME", label: "First Name", icon: Person },
  { value: "LAST_NAME", label: "Last Name", icon: Person },
  { value: "JOB_TITLE", label: "Job Title", icon: Work },
  { value: "COMPANY", label: "Company", icon: Business },
  { value: "CITY", label: "City", icon: LocationCity },
  { value: "COUNTRY", label: "Country", icon: Public },
  { value: "INDUSTRY", label: "Industry", icon: Business },
];

const CSV_REQUIRED_FIELDS = [
  "email",
  "firstName",
  "lastName",
  "jobTitle",
  "company",
  "city",
  "country",
  "industry",
];

const CSV_FIELD_MAP = {
  EMAIL: "email",
  FIRST_NAME: "firstName",
  LAST_NAME: "lastName",
  JOB_TITLE: "jobTitle",
  COMPANY: "company",
  CITY: "city",
  COUNTRY: "country",
  INDUSTRY: "industry",
};

function AddLeadModal({
  open,
  panelMode,
  emailListID,
  pageNum,
  onClose,
  setListState,
}) {
  const dispatch = useDispatch();
  const { notify } = useNotify();

  const [formData, setFormData] = useState(initialFields);
  const [csvFile, setCsvFile] = useState(null);
  const [dragActive, setDragActive] = useState(false);
  const [csvHeaders, setCsvHeaders] = useState([]);
  const [csvRows, setCsvRows] = useState([]);
  const [allCsvRows, setAllCsvRows] = useState([]);
  const [columnMapping, setColumnMapping] = useState({});
  const [isParsingCsv, setIsParsingCsv] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const resetState = () => {
    setFormData(initialFields);
    setCsvFile(null);
    setCsvHeaders([]);
    setCsvRows([]);
    setAllCsvRows([]);
    setColumnMapping({});
    setDragActive(false);
    setIsParsingCsv(false);
    setIsSubmitting(false);
  };

  useEffect(() => {
    if (!open) {
      resetState();
    }
  }, [open]);

  const parseCsvContent = (text) => {
    const rows = [];
    let currentRow = [];
    let currentCell = "";
    let inQuotes = false;

    for (const char of text) {
      if (char === '"') {
        inQuotes = !inQuotes;
        continue;
      }

      if (char === "," && !inQuotes) {
        currentRow.push(currentCell.trim());
        currentCell = "";
        continue;
      }

      if (char === "\n" && !inQuotes) {
        currentRow.push(currentCell.trim());
        if (currentRow.some((cell) => cell !== "")) {
          rows.push(currentRow);
        }
        currentRow = [];
        currentCell = "";
        continue;
      }

      currentCell += char;
    }

    if (currentCell || currentRow.length) {
      currentRow.push(currentCell.trim());
      if (currentRow.some((cell) => cell !== "")) {
        rows.push(currentRow);
      }
    }

    return rows;
  };

  const handleDrop = (e) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);

    if (e.dataTransfer.files && e.dataTransfer.files[0]) {
      handleFile(e.dataTransfer.files[0]);
    }
  };

  const handleDragOver = (e) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(true);
  };

  const handleDragLeave = (e) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);
  };

  const handleFileChange = (e) => {
    handleFile(e.target.files[0]);
  };

  const handleRemoveFile = (e) => {
    e.stopPropagation();
    setCsvFile(null);
    setCsvHeaders([]);
    setCsvRows([]);
    setAllCsvRows([]);
    setColumnMapping({});
    setIsParsingCsv(false);
    const input = document.getElementById("csvInput");
    if (input) input.value = "";
  };

  const handleFile = (file) => {
    if (!file) return;
    const isCSV =
      file.type === "text/csv" || file.name.toLowerCase().endsWith(".csv");
    if (!isCSV) {
      notify("Please upload a CSV file", "error");
      return;
    }

    setCsvFile(file);
    setCsvHeaders([]);
    setCsvRows([]);
    setAllCsvRows([]);
    setColumnMapping({});
    setIsParsingCsv(true);

    const reader = new FileReader();
    reader.onload = (event) => {
      const text = event.target.result || "";
      const parsedRows = parseCsvContent(text).filter((row) =>
        row.some((value) => value.trim() !== ""),
      );

      if (!parsedRows.length) {
        notify("The uploaded CSV file is empty or invalid.", "error");
        setCsvFile(null);
        setIsParsingCsv(false);
        return;
      }

      const headers = parsedRows[0].map((header) => header.trim());
      const dataRows = parsedRows.slice(1);

      if (!headers.some(Boolean) || dataRows.length === 0) {
        notify(
          "The uploaded CSV file does not contain any lead rows.",
          "error",
        );
        setCsvFile(null);
        setIsParsingCsv(false);
        return;
      }

      const mapping = {};
      let emailAlreadyMapped = false;
      headers.forEach((header, index) => {
        const normalized = (header || "")
          .toLowerCase()
          .replace(/[^a-z0-9]/g, "");

        if (normalized.includes("email")) {
          if (!emailAlreadyMapped) {
            mapping[index] = "EMAIL";
            emailAlreadyMapped = true;
          } else {
            mapping[index] = "IGNORE";
          }
        } else if (normalized.includes("firstname")) {
          mapping[index] = "FIRST_NAME";
        } else if (normalized.includes("lastname")) {
          mapping[index] = "LAST_NAME";
        } else if (
          normalized.includes("jobtitle") ||
          normalized.includes("job") ||
          normalized.includes("title")
        ) {
          mapping[index] = "JOB_TITLE";
        } else if (normalized.includes("company")) {
          mapping[index] = "COMPANY";
        } else if (normalized.includes("city")) {
          mapping[index] = "CITY";
        } else if (normalized.includes("country")) {
          mapping[index] = "COUNTRY";
        } else if (normalized.includes("industry")) {
          mapping[index] = "INDUSTRY";
        } else {
          mapping[index] = "IGNORE";
        }
      });

      setCsvHeaders(headers);
      setCsvRows(dataRows.slice(0, 4));
      setAllCsvRows(dataRows);
      setColumnMapping(mapping);
      setIsParsingCsv(false);
    };

    reader.onerror = () => {
      notify("Unable to read the selected CSV file.", "error");
      setCsvFile(null);
      setIsParsingCsv(false);
    };

    reader.readAsText(file);
  };

  const validateCsvMapping = () => {
    if (!csvHeaders.length) return "Please upload a valid CSV file.";
    if (!Object.values(columnMapping).includes("EMAIL")) {
      return "Please map at least one CSV column to Email.";
    }
    return "";
  };

  const buildProcessedCsvFile = () => {
    const rows = [CSV_REQUIRED_FIELDS];

    allCsvRows.forEach((row) => {
      const record = {
        email: "",
        firstName: "",
        lastName: "",
        jobTitle: "",
        company: "",
        city: "",
        country: "",
        industry: "",
      };

      csvHeaders.forEach((header, index) => {
        const key = CSV_FIELD_MAP[columnMapping[index]];
        if (key) {
          record[key] = (row[index] || "").trim();
        }
      });

      rows.push([
        record.email,
        record.firstName,
        record.lastName,
        record.jobTitle,
        record.company,
        record.city,
        record.country,
        record.industry,
      ]);
    });

    const csvContent = rows
      .map((row) =>
        row
          .map((cell) => `"${String(cell ?? "").replace(/"/g, '""')}"`)
          .join(","),
      )
      .join("\n");

    return new File([csvContent], csvFile.name, {
      type: "text/csv",
    });
  };

  const handleChange = (field, value) => {
    setFormData((prev) => ({
      ...prev,
      [field]: { value, error: "" },
    }));
  };

  const handleOnBlur = (name, value) => {
    const trimmedValue = typeof value === "string" ? value.trim() : value;
    const error = getValidationError(name, trimmedValue);
    setFormData((prev) => ({
      ...prev,
      [name]: {
        ...prev[name],
        value: trimmedValue,
        error,
      },
    }));
  };

  const getValidationError = (name, value) => {
    const trimmedValue = typeof value === "string" ? value.trim() : value;

    switch (name) {
      case "email":
        if (formData.type.value !== "single") return "";
        return !trimmedValue
          ? "Please enter email"
          : !REGEX.EMAIL_REGEX.test(trimmedValue)
          ? "Enter valid email"
          : "";
      case "csvFile":
        if (formData.type.value !== "csv") return "";
        return !csvFile ? "Please upload a CSV file" : "";
      default:
        return "";
    }
  };

  const handleTypeChange = (value) => {
    setFormData({
      ...initialFields,
      type: { value, error: "" },
    });

    setCsvFile(null);
    setCsvHeaders([]);
    setCsvRows([]);
    setAllCsvRows([]);
    setColumnMapping({});
    setIsParsingCsv(false);

    const input = document.getElementById("csvInput");
    if (input) input.value = "";
  };

  const validateForm = () => {
    const newData = { ...formData };
    Object.keys(newData).forEach((key) => {
      if (key === "id") return;
      newData[key] = {
        ...newData[key],
        error: getValidationError(key, newData[key].value),
      };
    });
    setFormData(newData);
    return Object.keys(newData)
      .filter((key) => key !== "id")
      .every((key) => newData[key].error === "");
  };

  // const startVerificationPolling = () => {
  //   if (window.leadVerificationInterval) {
  //     clearInterval(window.leadVerificationInterval);
  //   }

  //   setVerificationState({
  //     isChecking: true,
  //     pendingCount: 0,
  //   });

  //   window.leadVerificationInterval = setInterval(() => {
  //     getLeadsList(dispatch, emailListID, pageNum);
  //   }, 3000);
  // };

  const handleSubmit = async () => {
    if (!validateForm()) return;

    if (formData.type.value === "csv" && !csvFile) {
      notify("Please upload a CSV file", "error");
      return;
    }

    try {
      setIsSubmitting(true);
      setListState((prev) => ({ ...prev, showLoaderOnClick: true }));

      const payload = new FormData();
      payload.append("type", formData.type.value);

      if (formData.type.value === "single") {
        payload.append("email", formData.email.value.trim());
        payload.append("firstName", formData.firstName.value.trim());
        payload.append("lastName", formData.lastName.value.trim());
        payload.append("jobTitle", formData.jobTitle.value.trim());
        payload.append("company", formData.company.value.trim());
        // payload.append("country", formData.country.value);
        // payload.append("city", formData.city.value);
        // payload.append("industry", formData.industry.value);
      } else {
        const mappingError = validateCsvMapping();
        if (mappingError) {
          notify(mappingError, "error");
          return;
        }

        const processedFile = buildProcessedCsvFile();
        payload.append("file", processedFile);
      }

      payload.append("emailListID", String(emailListID));
      const res = await addLeadList(payload);

      if (res?.code === 200) {
        notify(res?.message || "Lead added successfully", "success");
        resetState();
        onClose();
        await getLeadsList(dispatch, emailListID, pageNum);
        // startVerificationPolling();
      } else {
        notify(res?.message || "Failed to add lead", "error");
      }
    } catch (error) {
      notify(error?.message || "Something went wrong", "error");
    } finally {
      setIsSubmitting(false);
      setListState((prev) => ({ ...prev, showLoaderOnClick: false }));
    }
  };

  const renderMappingValue = (value) => {
    const option = CSV_MAPPING_FIELDS.find((item) => item.value === value);
    if (!option) return "";
    const IconComponent = option.icon;
    return (
      <MDBox display="flex" alignItems="center" gap={1}>
        <IconComponent
          fontSize="small"
          sx={{
            color:
              value === "IGNORE" ? colors.vividRed.main : colors.deepBlue.main,
          }}
        />
        <MDTypography
          variant="button"
          sx={{ fontWeight: 400, fontSize: "0.75rem" }}
        >
          {option.label}
        </MDTypography>
      </MDBox>
    );
  };

  const getAvailableOptions = (currentIndex) => {
    return CSV_MAPPING_FIELDS.map((option) => {
      const isUsedElsewhere = Object.entries(columnMapping).some(
        ([index, value]) =>
          Number(index) !== currentIndex && value === option.value,
      );

      return {
        ...option,
        isDisabled: option.value !== "IGNORE" && isUsedElsewhere,
      };
    });
  };

  return (
    <CenterPopup isOpen={open} onClose={onClose} width="900px">
      <MDBox sx={{ p: 2 }}>
        <MDTypography variant="h6" sx={leadListUiStyles.modalTitle}>
          {panelMode === "edit" ? "Edit Email List" : "Add Email List"}
        </MDTypography>

        <Grid container spacing={2}>
          <Grid item xs={12}>
            <RequiredLabel>Select Type</RequiredLabel>
            <RadioGroup
              row
              name="type"
              value={formData.type.value}
              onChange={(e) => handleTypeChange(e.target.value)}
            >
              <FormControlLabel
                value="single"
                control={<Radio size="small" />}
                label="Single Email"
                sx={radioLabelSx}
              />
              <FormControlLabel
                value="csv"
                control={<Radio size="small" />}
                label="CSV Upload"
                sx={radioLabelSx}
              />
            </RadioGroup>
          </Grid>

          {formData.type.value === "single" ? (
            <>
              <Grid item xs={6}>
                <RequiredLabel required>Email</RequiredLabel>
                <MDInput
                  variant="outlined"
                  type="email"
                  name="email"
                  placeholder="Enter email"
                  value={formData.email.value}
                  onChange={(e) => handleChange(e.target.name, e.target.value)}
                  onBlur={(e) => handleOnBlur(e.target.name, e.target.value)}
                  error={!!formData.email.error}
                  fullWidth
                />
                {formData.email.error && (
                  <MDTypography variant="caption" color="error">
                    {formData.email.error}
                  </MDTypography>
                )}
              </Grid>

              <Grid item xs={6}>
                <RequiredLabel>First Name</RequiredLabel>
                <MDInput
                  type="text"
                  name="firstName"
                  placeholder="Enter first name"
                  value={formData.firstName.value}
                  onChange={(e) => handleChange(e.target.name, e.target.value)}
                  inputProps={{ maxLength: 50 }}
                  fullWidth
                />
              </Grid>

              <Grid item xs={6}>
                <RequiredLabel>Last Name</RequiredLabel>
                <MDInput
                  name="lastName"
                  placeholder="Enter last name"
                  value={formData.lastName.value}
                  onChange={(e) => handleChange(e.target.name, e.target.value)}
                  inputProps={{ maxLength: 50 }}
                  fullWidth
                />
              </Grid>

              <Grid item xs={6}>
                <RequiredLabel>Job Title</RequiredLabel>
                <MDInput
                  name="jobTitle"
                  placeholder="Enter job title"
                  value={formData.jobTitle.value}
                  onChange={(e) => handleChange(e.target.name, e.target.value)}
                  inputProps={{ maxLength: 100 }}
                  fullWidth
                />
              </Grid>

              <Grid item xs={6}>
                <RequiredLabel>Company</RequiredLabel>
                <MDInput
                  name="company"
                  placeholder="Enter company"
                  value={formData.company.value}
                  onChange={(e) => handleChange(e.target.name, e.target.value)}
                  inputProps={{ maxLength: 100 }}
                  fullWidth
                />
              </Grid>

              {/* <Grid item xs={6}>
                <RequiredLabel>Country</RequiredLabel>
                <Select
                  fullWidth
                  name="country"
                  placeholder="Select country"
                  value={formData.country.value}
                  IconComponent={ArrowDropDown}
                  onChange={(e) => handleChange("country", e.target.value)}
                  sx={selectSx("100%")}
                />
              </Grid>

              <Grid item xs={6}>
                <RequiredLabel>City</RequiredLabel>
                <Select
                  fullWidth
                  name="city"
                  placeholder="Select city"
                  value={formData.city.value}
                  IconComponent={ArrowDropDown}
                  onChange={(e) => handleChange("city", e.target.value)}
                  sx={selectSx("100%")}
                />
              </Grid> */}
            </>
          ) : (
            <Grid item xs={12}>
              <RequiredLabel required>Upload CSV</RequiredLabel>

              <MDBox
                onDrop={handleDrop}
                onDragOver={handleDragOver}
                onDragLeave={handleDragLeave}
                onClick={() =>
                  !csvFile && document.getElementById("csvInput").click()
                }
                sx={{
                  ...leadCsvStyles.uploadBox,
                  borderColor: dragActive
                    ? colors.deepBlue.main
                    : leadCsvStyles.uploadBox.borderColor,
                  backgroundColor: dragActive
                    ? colors.grey[100]
                    : csvFile
                    ? colors.white.main
                    : leadCsvStyles.uploadBox.backgroundColor,
                  cursor: csvFile ? "default" : "pointer",
                }}
              >
                <input
                  id="csvInput"
                  type="file"
                  accept=".csv"
                  hidden
                  onChange={handleFileChange}
                />

                {csvFile && (
                  <MDBox
                    onClick={handleRemoveFile}
                    sx={leadCsvStyles.uploadRemoveButton}
                  >
                    <i
                      className="fas fa-times"
                      style={{ color: colors.white.main, fontSize: "0.75rem" }}
                    />
                  </MDBox>
                )}

                {!csvFile ? (
                  <>
                    <Icon sx={{ fontSize: 45, color: "primary.main" }}>
                      cloud_upload
                    </Icon>
                    <MDTypography
                      variant="body2"
                      sx={leadListUiStyles.uploadSubtitle}
                    >
                      Drag & drop your CSV file here or click to browse
                    </MDTypography>
                  </>
                ) : (
                  <MDBox
                    display="flex"
                    flexDirection="column"
                    alignItems="center"
                  >
                    <MDTypography
                      variant="body2"
                      sx={{ ...leadListUiStyles.uploadTitle, mb: 1 }}
                    >
                      {csvFile.name}
                    </MDTypography>
                    <MDTypography
                      variant="caption"
                      sx={leadListUiStyles.helperText}
                    >
                      {isParsingCsv ? "Parsing CSV..." : "CSV File uploaded"}
                    </MDTypography>
                  </MDBox>
                )}
              </MDBox>

              {isParsingCsv && (
                <MDBox
                  mt={2}
                  display="flex"
                  alignItems="center"
                  justifyContent="center"
                >
                  <i className="fa fa-spinner fa-spin" />
                  <MDTypography
                    variant="caption"
                    sx={{ ...leadListUiStyles.helperText, ml: 1 }}
                  >
                    Parsing CSV columns...
                  </MDTypography>
                </MDBox>
              )}

              {csvHeaders.length > 0 && !isParsingCsv && (
                <MDBox mt={2}>
                  <MDBox sx={leadCsvStyles.mappingContainer}>
                    <MDBox px={2} py={1} sx={leadCsvStyles.mappingHeader} />
                    <MDBox sx={{ overflowX: "auto" }}>
                      <table
                        style={{ width: "100%", borderCollapse: "collapse" }}
                      >
                        <thead>
                          <tr style={{ backgroundColor: colors.grey[100] }}>
                            <th style={leadCsvStyles.tableHeaderCell}>
                              CSV Column
                            </th>
                            <th style={leadCsvStyles.tableHeaderCell}>
                              Select Type
                            </th>
                            <th style={leadCsvStyles.tableHeaderCell}>
                              Sample Records
                            </th>
                          </tr>
                        </thead>
                        <tbody>
                          {csvHeaders.map((header, index) => (
                            <tr
                              key={`${header}-${index}`}
                              style={leadCsvStyles.tableRow}
                            >
                              <td style={leadCsvStyles.tableCell}>
                                <MDTypography variant="caption">
                                  {header?.trim() ? header : "Empty Column"}
                                </MDTypography>
                              </td>
                              <td
                                style={{
                                  ...leadCsvStyles.tableCell,
                                  minWidth: "220px",
                                }}
                              >
                                <Select
                                  fullWidth
                                  size="small"
                                  value={columnMapping[index] || "IGNORE"}
                                  sx={{
                                    ...selectSx("100%"),
                                    "& .MuiOutlinedInput-root": {
                                      minHeight: "2.5rem",
                                    },
                                    "& .MuiSelect-select": {
                                      display: "flex",
                                      alignItems: "center",
                                      py: "0.6rem",
                                    },
                                  }}
                                  renderValue={renderMappingValue}
                                  onChange={(e) =>
                                    setColumnMapping((prev) => ({
                                      ...prev,
                                      [index]: e.target.value,
                                    }))
                                  }
                                >
                                  {getAvailableOptions(index).map((option) => {
                                    const IconComponent = option.icon;
                                    return (
                                      <MenuItem
                                        key={option.value}
                                        value={option.value}
                                        disabled={option.isDisabled}
                                        sx={{
                                          minHeight: "2.4rem",
                                          py: 0.75,
                                          "&.Mui-disabled": { opacity: 0.45 },
                                        }}
                                      >
                                        <ListItemIcon
                                          sx={{
                                            minWidth: 34,
                                            color: option.isDisabled
                                              ? colors.grey[400]
                                              : option.value === "IGNORE"
                                              ? colors.vividRed.main
                                              : colors.deepBlue.main,
                                          }}
                                        >
                                          <IconComponent fontSize="small" />
                                        </ListItemIcon>
                                        <ListItemText
                                          primary={option.label}
                                          primaryTypographyProps={{
                                            fontSize: "0.85rem",
                                            fontWeight: 500,
                                            color: option.isDisabled
                                              ? colors.grey[500]
                                              : colors.dark.main,
                                          }}
                                        />
                                      </MenuItem>
                                    );
                                  })}
                                </Select>
                              </td>
                              <td style={leadCsvStyles.tableCell}>
                                <MDBox
                                  display="flex"
                                  flexDirection="column"
                                  gap={0.5}
                                >
                                  {csvRows.map((row, rowIndex) => (
                                    <MDTypography
                                      key={`${header}-${index}-${rowIndex}`}
                                      variant="caption"
                                      sx={leadCsvStyles.sampleText}
                                    >
                                      {row[index] || "-"}
                                    </MDTypography>
                                  ))}
                                </MDBox>
                              </td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </MDBox>
                  </MDBox>
                </MDBox>
              )}
            </Grid>
          )}

          <Grid item xs={12}>
            <MDBox display="flex" justifyContent="flex-end">
              <MDButton
                variant="gradient"
                color="light"
                onClick={onClose}
                sx={leadListUiStyles.secondaryButton}
              >
                Cancel
              </MDButton>
              <MDButton
                type="submit"
                variant="contained"
                color="deepBlue"
                onClick={handleSubmit}
                disabled={isSubmitting}
              >
                {isSubmitting ? (
                  <i className="fa fa-spinner fa-spin" />
                ) : (
                  "Proceed"
                )}
              </MDButton>
            </MDBox>
          </Grid>
        </Grid>
      </MDBox>
    </CenterPopup>
  );
}

export default AddLeadModal;
