import { useEffect, useState } from 'react';
import { Outlet, useLocation, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthProvider';
import React from 'react';
import useAuthStore from '../stores/auth';

const ProtectedRoute = () => {
  const { setUser, setAuthorization, logout } = useAuthStore();
  const userPromise = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  const [authenticated, setAuthenticated] = useState(false);

  const denyAccessAndRedirect = () => {
    logout();
    setAuthenticated(false);
    navigate('/signup', { replace: true, state: { from: location } });
  };

  useEffect(() => {
    if (!userPromise) return denyAccessAndRedirect();
    userPromise
      .then((user) => {
        setAuthenticated(user !== null);

        setUser(user);
        setAuthorization();
      })
      .catch(() => denyAccessAndRedirect());
  }, [userPromise]);

  return authenticated !== null ? <Outlet /> : null;
};

export default ProtectedRoute;
