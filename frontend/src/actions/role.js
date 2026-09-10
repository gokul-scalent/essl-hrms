import { roleList } from "constants/role";
import * as API from "../API/index";
import { authorizedUser } from "components/CommonComponent/CommonFunction";

//get role list
export const getRoleListing = async (dispatch) => {
  try {
    const res = await API.get_role_list();
    dispatch({ type: roleList, payload: res?.data });
  } catch (error) {
    authorizedUser(error.response?.data);
    dispatch({type: roleList,payload: error.response?.data ?? { message: error?.message },
    });
  }
};