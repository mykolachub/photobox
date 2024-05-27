import React from 'react';
import authStore from '../../stores/auth';
import { useNavigate } from 'react-router-dom';
import AsyncImage from '../AsyncImage';

import './Header.css';
import ButtonSmall from '../Buttons/ButtonSmall';

const Header = () => {
  const { user, logout } = authStore();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/signup');
  };

  return (
    <div className="header__wrapper">
      <header className="header">
        <div className="header__logo header__item">Photobox</div>

        <nav className="header__nav header__item">
          <div>Home</div>
        </nav>

        <div className="header__profile header__item">
          {user.email && (
            <>
              <div className="profile">
                <div className="profile__picture">
                  <AsyncImage src={user.picture} />
                </div>
                {user.username}
              </div>
              <ButtonSmall message="Logout" onClick={handleLogout} />
            </>
          )}
          {/* {!user.email && (
            <Link to="/signup" style={{ textDecoration: 'none' }}>
              <ButtonSmall message="Signup" onClick={handleLogout} />
            </Link>
          )} */}
        </div>
      </header>
    </div>
  );
};

export default Header;
