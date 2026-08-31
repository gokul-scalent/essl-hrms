import { employeeList } from "constants/employee";

const initialState = {
  employeeList: {},
};

const employee = (state = initialState, action) => {
  switch (action.type) {
    case employeeList:
      return { ...state, employeeList: action.payload };

    default:
      return state;
  }
};
export default employee;
