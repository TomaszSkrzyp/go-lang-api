import React, { createContext, useState, useEffect } from "react";
import PropTypes from 'prop-types';
export const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const [token, setToken] = useState(() => localStorage.getItem("token") || null);

  useEffect(() => {
    if (token) localStorage.setItem("token", token);
    else localStorage.removeItem("token");
  }, [token]);
  const logout = () => {
    setToken(null);
  };

  return (
    <AuthContext.Provider value={{ token,setToken, logout }}>
      {children}
    </AuthContext.Provider>
  );
}
AuthProvider.propTypes = {
  children: PropTypes.node.isRequired,
};