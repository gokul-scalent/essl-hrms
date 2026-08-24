import * as api from "../API/index";

// Login API
export const login = async (loginData) => {
  try {
    const res = await api.login_api(loginData);
    return res.data;
  } catch (error) {
    return error?.response?.data || error;
  }
};

//logout api
export const logout = async () => {
  try {
    const res = await api.logout_api();
    return res.data;
  } catch (error) {
    return error?.response?.data || error;
  }
};
