import { userList } from "constants/users";
import * as API from "../API/index";
import { authorizedUser } from "components/CommonComponent/CommonFunction";

//get users list
export const getUsersList = async (dispatch,pageNum,filters,searchString,) => {
  try {
    const res = await API.get_users_list(pageNum,filters,searchString,);
    dispatch({ type: userList, payload: res.data });
  } catch (error) {
    authorizedUser(error.response?.data);
    dispatch({ type: userList, payload: error?.response?.data || error });
  }
};

export const addUser = async (bodyData) => {
  try {
    const res = await API.add_users(bodyData);
    return res?.data;
  } catch (error) {
    authorizedUser(error.response?.data);
    return error?.response?.data || error;
  }
};

export const getUserDetailById = async (id) => {
  try {
    const res = await API.get_user_details_by_user(id);
    return res?.data;
  } catch (error) {
    return error?.response?.data || error;
  }
};

export const updateUserDetail = async (id, bodyData) => {
  try {
    const res = await API.update_user_by_id(id, bodyData);
    return res?.data;
  } catch (error) {
    return error?.response?.data || error;
  }
};

export const deleteUser = async (id) => {
  try {
    const res = await API.delete_user_by_ID(id);
    return res?.data;
  } catch (error) {
    return error?.response?.data || error;
  }
};