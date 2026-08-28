import { employeeList } from "constants/employee";
import * as API from "../API/index";
import { authorizedUser } from "components/CommonComponent/CommonFunction";

//get users list
export const getEmployeeList = async (dispatch,pageNum,filters,searchString,) => {
  try {
    const res = await API.get_employee_list(pageNum,filters,searchString,);
    dispatch({ type: employeeList, payload: res.data });
  } catch (error) {
    authorizedUser(error.response?.data);
    dispatch({ type: employeeList, payload: error?.response?.data || error });
  }
};