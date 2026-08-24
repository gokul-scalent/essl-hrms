import { emailList } from "constants/email";
import * as API from "../API/index";
import { authorizedUser } from "components/CommonComponent/CommonFunction";

//get list
export const getEmailLists = async (dispatch,pageNum,filters,searchString,) => {
  try {
    const res = await API.email_list(pageNum, filters, searchString);
    dispatch({ type: emailList, payload: res.data });
  } catch (error) {
    authorizedUser(error.response?.data);
    dispatch({ type: emailList, payload: error?.response?.data || error });
  }
};

//add email list
export const addEmailList = async (bodyData) => {
  try {
    const res = await API.add_email_list(bodyData);
    return res?.data;
  } catch (error) {
    authorizedUser(error.response?.data);
    return error?.response?.data || error;
  }
};

//get email list by id
export const getEmailListById = async (id) => {
  try {
    const res = await API.get_email_list_by_ID(id);
    return res?.data;
  } catch (error) {
    return error?.response?.data || error;
  }
};

//update email list by id
export const updateEmailListById = async (id, bodyData) => {
  try {
    const res = await API.edit_email_list(id, bodyData);
    return res?.data;
  } catch (error) {
    return error?.response?.data || error;
  }
};

export const deleteEmailListById = async (id) => {
  try {
    const res = await API.delete_email_list_by_ID(id);
    return res?.data;
  } catch (error) {
    return error?.response?.data || error;
  }
};
