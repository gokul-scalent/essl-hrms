import axios from "axios";
import {
  BASE_URL,
  LOGIN,
  LOGOUT,
  USERS_LIST,
  USERS,
  CHANGE_PASSWORD,
  EMPLOYEE_LIST,
  ATTENDANCE_LOG_LIST,
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

//get the employee list
export const get_employee_list = (pageNum, filtersJSON, searchString) =>
  API.get(EMPLOYEE_LIST, {
    params: {
      page: pageNum,
      ...(filtersJSON ? { filtersJSON } : {}),
      ...(searchString ? { searchString } : {}),
    },
  })

//get the attendance log list
export const get_attendance_log_list = (pageNum, filtersJSON, searchString) =>
  API.get(ATTENDANCE_LOG_LIST, {
    params: {
      page: pageNum,
      ...(filtersJSON ? { filtersJSON } : {}),
      ...(searchString ? { searchString } : {}),
    },
  })