import React from 'react';
import authStore from '../stores/auth';

import './Signup.css';
import GoogleIcon from '../assets/images/icon-google.svg';

const Signup = () => {
  const { signupGoogle } = authStore();

  const handleSignupGoogle = async () => {
    const { url } = await signupGoogle();
    window.open(url, '_self');
  };

  return (
    <div className="signup__wrapper">
      <div className="signup__content">
        <span className="signup__logo">Photobox</span>
        <button onClick={handleSignupGoogle} className="signup__google">
          <img src={GoogleIcon} alt="Google Icon" />
          Signup via Google
        </button>
      </div>
    </div>
  );
};

export default Signup;
