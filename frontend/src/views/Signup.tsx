import React from 'react';
import authStore from '../stores/auth';

const Signup = () => {
  const { signupGoogle } = authStore();

  const handleSignupGoogle = async () => {
    const { url } = await signupGoogle();
    window.open(url, '_self');
  };

  return (
    <div>
      <button onClick={handleSignupGoogle}>signup via google</button>
    </div>
  );
};

export default Signup;
