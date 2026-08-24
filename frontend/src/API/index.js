import axios from "axios";
import {
  BASE_URL,
  EMAIL_lIST_LISTING,
  LOGIN,
  EMAIL_LIST,
  LEADS_LISTING,
  LEADS,
  LOGOUT,
  REVERIFY_LEAD,
  SETTING_CONFIGURATION,
  SETTING_CONFIGURATION_LIST,
  USERS_LIST,
  USERS,
  CHANGE_PASSWORD,
} from "./apiConstants";

const API = axios.create({
  //baseURL: BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

API.interceptors.request.use((req) => {
  if (localStorage.getItem("p")) {
    req.headers.token = JSON.parse(localStorage.getItem("p")).token;
  }
  return req;
});

// Login API
export const login_api = (loginData) => API.post(`${LOGIN}`, loginData);

// Logout API
export const logout_api = () => API.post(`${LOGOUT}`);

//email list api
export const email_list = (pageNum, filtersJSON, searchString) =>
  API.get(EMAIL_lIST_LISTING, {
    params: {
      page: pageNum,
      ...(filtersJSON ? { filtersJSON } : {}),
      ...(searchString ? { searchString } : {}),
    },
  });

export const add_email_list = (bodyData) => API.post(EMAIL_LIST, bodyData);
export const edit_email_list = (id, bodyData) =>
  API.patch(`${EMAIL_LIST}${id}`, bodyData);
export const get_email_list_by_ID = (id) => API.get(`${EMAIL_LIST}${id}`);
export const delete_email_list_by_ID = (id) => API.delete(`${EMAIL_LIST}${id}`);

//leads api
export const leads_list = (emailListID, pageNum, filtersJSON, searchString) =>
  API.get(LEADS_LISTING, {
    params: {
      emailListID,
      page: pageNum,
      ...(filtersJSON ? { filtersJSON } : {}),
      ...(searchString ? { searchString } : {}),
    },
  });

// export const add_leads = (formData) => API.post(LEADS, formData);
export const add_leads = (formData) =>
  API.post(LEADS, formData, {
    headers: {
      "Content-Type": "multipart/form-data",
    },
  });
export const edit_leads = (id, bodyData) =>
  API.patch(`${LEADS}${id}`, bodyData);
export const get_leads_by_ID = (id) => API.get(`${LEADS}${id}`);
export const delete_leads_by_ID = (id) => API.delete(`${LEADS}${id}`);

//download safe csv
export const safe_download_csv = (params) =>
  API.get(`${LEADS}download`, {
    params,
    responseType: "blob",
  });

//reverify the lead
export const reverify_lead = (id) => API.put(`${REVERIFY_LEAD}${id}`);

//setting configuration
export const add_setting_configuration =(bodyData) => API.post(SETTING_CONFIGURATION, bodyData);
export const update_setting_configuration = (id, bodyData) =>API.patch(`${SETTING_CONFIGURATION}${id}`,bodyData);

//list setting configuration
export const get_user_setting_list = () => API.get(SETTING_CONFIGURATION_LIST);

//get the users list
export const get_users_list = (pageNum, filtersJSON, searchString) =>
  API.get(USERS_LIST, {
    params: {
      page: pageNum,
      ...(filtersJSON ? { filtersJSON } : {}),
      ...(searchString ? { searchString } : {}),
    },
  });

//add users
export const add_users = (bodyData) => API.post(`${USERS}`,bodyData);
//update user list
export const update_user_by_id =(id, bodyData) => API.patch(`${USERS}${id}`, bodyData);
//get user details by id
export const get_user_details_by_user =(id) => API.get(`${USERS}${id}`);
export const delete_user_by_ID = (id) => API.delete(`${USERS}${id}`);

//change password
export const change_password = (passwordData) =>
  API.patch(`${CHANGE_PASSWORD}`, passwordData);
