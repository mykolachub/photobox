import React from 'react';
import { Outlet, useNavigate } from 'react-router-dom';
import authStore from '../stores/auth';
import AsyncImage from './AsyncImage';

const Layout = () => {
  const { user, logout } = authStore();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/signup');
  };

  return (
    <div className="wrapper">
      <header>
        header
        {user.email && (
          <>
            <p>
              Welcome {user.email} {user.username}
            </p>
            <AsyncImage src={user.picture} />
            <button onClick={handleLogout}>Logout</button>
          </>
        )}
      </header>
      <main>
        <Outlet />
      </main>
      <footer className="footer">Photobox 2024</footer>
    </div>
  );
};

export default Layout;
