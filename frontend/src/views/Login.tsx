import React, { useEffect, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import authStore from '../stores/auth';

const Login = () => {
  const { loginGoogle, setAuthorization } = authStore();

  const navigate = useNavigate();
  const location = useLocation();
  const from = location.state?.from?.pathname || '/';

  const [code] = useState<string>(() => {
    const queryParams = new URLSearchParams(window.location.search);
    return queryParams.get('code') || '';
  });

  useEffect(() => {
    loginGoogle({ code })
      .then(() => {
        setAuthorization();
        navigate(from, { replace: true });
      })
      .catch((err) => {
        console.log(err);
      });
  }, []);

  return <div>processing...</div>;
};

export default Login;
