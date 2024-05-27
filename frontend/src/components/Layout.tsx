import React from 'react';
import { Outlet } from 'react-router-dom';
import Header from './Header/Header';

const Layout = () => {
  return (
    <div className="wrapper">
      <Header />
      <main>
        <Outlet />
      </main>
      <footer className="footer">Photobox 2024</footer>
    </div>
  );
};

export default Layout;
