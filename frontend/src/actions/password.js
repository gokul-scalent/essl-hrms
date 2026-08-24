import * as API from "../API/index";

//change password
export const changePasswordAPI = async (passwordData) => {
  try {
    const res = await API.change_password(passwordData);
    return res.data;
  } catch (err) {
    return err?.response?.data;
  }
};