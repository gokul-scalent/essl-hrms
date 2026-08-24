import React, { useState } from "react";
import { createRoot } from "react-dom/client";
import SweetAlert from "react-bootstrap-sweetalert";
import colors from "assets/theme/base/colors";

export const showAlert = ({
  title = "",
  message = "",
  type = "info",
  showCancel = false,
  confirmText = "OK",
  cancelText = "Cancel",
  onConfirm = () => {},
  onCancel = () => {},
  confirmColor = colors.deepBlue.main,
  cancelColor = colors.light.main,
  buttonGap = "0.625rem",
  requireInput = false,
  inputLabel = "Remark",
}) => {
  const container = document.createElement("div");
  const root = createRoot(container);
  document.body.appendChild(container);

  const cleanup = () => {
    root.unmount();
    container.remove();
  };

  const AlertComponent = () => {
    const [inputValue, setInputValue] = useState("");
    const [error, setError] = useState("");

    const handleConfirm = () => {
      if (requireInput && !inputValue.trim()) {
        setError(`Please enter ${inputLabel}`);
        return;
      }
      cleanup();
      onConfirm(inputValue);
    };

    const handleCancel = () => {
      cleanup();
      onCancel();
    };

    let alertTypeProps = {};
    switch (type) {
      case "success":
        alertTypeProps = { success: true };
        break;
      case "error":
        alertTypeProps = { danger: true };
        break;
      case "warning":
        alertTypeProps = { warning: true };
        break;
      case "info":
      default:
        alertTypeProps = { info: true };
        break;
    }

    return (
      <SweetAlert
        {...alertTypeProps}
        showCancel={showCancel}
        title={
          <span
            style={{
              fontSize: "1.125rem",
              lineHeight: "1.3rem",
              display: "inline-block",
              whiteSpace: "normal",
            }}
          >
            {title}
          </span>
        }
        onConfirm={handleConfirm}
        onCancel={handleCancel}
        style={{ width: "28.125rem" }}
        customButtons={
          <div
            style={{
              display: "flex",
              justifyContent: "center",
              gap: buttonGap,
            }}
          >
            {showCancel && (
              <button
                onClick={handleCancel}
                style={{
                  backgroundColor: cancelColor,
                  color: "#000",
                  border: "none",
                  padding: "0.625rem 1rem",
                  borderRadius: "0.25rem",
                  cursor: "pointer",
                }}
              >
                {cancelText}
              </button>
            )}

            <button
              onClick={handleConfirm}
              style={{
                backgroundColor: confirmColor,
                color: "#fff",
                border: "none",
                padding: "0.625rem 1rem",
                borderRadius: "0.25rem",
                cursor: "pointer",
              }}
            >
              {confirmText}
            </button>
          </div>
        }
      >
        <div style={{ fontSize: "1rem" }}>
          {message}

          {requireInput && (
            <div style={{ textAlign: "left", marginTop: "15px" }}>
              <label
                style={{
                  fontSize: "14px",
                  fontWeight: 500,
                  display: "block",
                  marginBottom: "6px",
                }}
              >
                {inputLabel} <span style={{ color: "#d32f2f" }}>*</span>
              </label>

              <textarea
                value={inputValue}
                onChange={(e) => {
                  setInputValue(e.target.value);
                  if (error) setError("");
                }}
                placeholder={`Enter ${inputLabel}...`}
                rows={3}
                style={{
                  width: "100%",
                  padding: "8px 10px",
                  borderRadius: "6px",
                  border: error ? "1px solid #d32f2f" : "1px solid #ced4da",
                  fontSize: "14px",
                  outline: "none",
                }}
              />

              {error && (
                <p
                  style={{
                    color: "#d32f2f",
                    fontSize: "12px",
                    marginTop: "4px",
                  }}
                >
                  {error}
                </p>
              )}
            </div>
          )}
        </div>
      </SweetAlert>
    );
  };

  root.render(<AlertComponent />);
};
