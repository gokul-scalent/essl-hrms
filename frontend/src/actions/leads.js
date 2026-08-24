import { leadList } from "constants/leads";
import * as API from "../API/index";
import { authorizedUser } from "components/CommonComponent/CommonFunction";

//get list
export const getLeadsList = async (
  dispatch,
  emailListID,
  pageNum,
  filters,
  searchString,
) => {
  try {
    const res = await API.leads_list(
      emailListID,
      pageNum,
      filters,
      searchString,
    );
    dispatch({ type: leadList, payload: res.data });
  } catch (error) {
    authorizedUser(error.response?.data);
    dispatch({ type: leadList, payload: error?.response?.data || error });
  }
};

//add lead list
export const addLeadList = async (bodyData) => {
  try {
    const res = await API.add_leads(bodyData);
    return res?.data;
  } catch (error) {
    authorizedUser(error.response?.data);
    return error?.response?.data || error;
  }
};

//get lead list by id
export const getLeadByID = async (id) => {
  try {
    const res = await API.get_leads_by_ID(id);
    return res?.data;
  } catch (error) {
    return error?.response?.data || error;
  }
};

//update lead list by id
export const updateLeadListByID = async (id, bodyData) => {
  try {
    const res = await API.edit_leads(id, bodyData);
    return res?.data;
  } catch (error) {
    return error?.response?.data || error;
  }
};

export const deleteLeadByID = async (id) => {
  try {
    const res = await API.delete_leads_by_ID(id);
    return res?.data;
  } catch (error) {
    return error?.response?.data || error;
  }
};

export const DownloadSafeCsv = async (params) => {
  try {
    const res = await API.safe_download_csv(params);
    return res;
  } catch (error) {
    authorizedUser(error.response?.data);
    return error?.response?.data || error;
  }
};

//reverify the lead
export const ReverifyLeadById = async (id) => {
  try {
    const res = await API.reverify_lead(id);
    return res?.data;
  } catch (error) {
    authorizedUser(error.response?.data);
    return error?.response?.data || error;
  }
};
