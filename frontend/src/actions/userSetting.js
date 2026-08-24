import { authorizedUser } from "components/CommonComponent/CommonFunction";
import * as API from "../API/index";

//create user setting configuration
export const saveUserSettingConfiguration = async (bodyData) => {
  try {
    const res = await API.add_setting_configuration(bodyData);
    return res?.data;
  } catch (error) {
    authorizedUser(error.response?.data);
    return error?.response?.data || error;
  }
};

//update the user setting configuration
export const updateUserSettingConfiguration = async (id, bodyData) => {
  try {
    const res = await API.update_setting_configuration(id, bodyData);
    return res?.data;
  } catch (error) {
    authorizedUser(error.response?.data);
    return error?.response?.data || error;
  }
};

export const getUserSettingList = async () => {
  try {
    const res = await API.get_user_setting_list();
    return res?.data;
  } catch (error) {
    authorizedUser(error.response?.data);
    return error?.response?.data || error;
  }
};
