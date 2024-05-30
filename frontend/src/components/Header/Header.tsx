import React from 'react';
import authStore from '../../stores/auth';
import { useNavigate } from 'react-router-dom';
import AsyncImage from '../AsyncImage';

import './Header.css';
import ButtonSmall from '../Buttons/ButtonSmall';
import searchStore from '../../stores/search';

function bytesToGB(bytes: number) {
  const gigabytes = bytes / (1024 * 1024 * 1024) || 0;
  return gigabytes.toFixed(1);
}

function bytesToMB(bytes: number) {
  const gigabytes = bytes / (1024 * 1024) || 0;
  return gigabytes.toFixed(2);
}

const Header = () => {
  const { user, logout } = authStore();
  const navigate = useNavigate();
  const { setSearch } = searchStore();

  const handleLogout = () => {
    logout();
    navigate('/signup');
  };

  return (
    <div className="header__wrapper">
      <header className="header">
        <div className="header__logo header__item">Photobox</div>

        <div className="header__search header__item">
          <input
            type="search"
            name="search"
            id="search"
            onChange={(e) => {
              setSearch(e.target.value.toLowerCase());
            }}
            className="header__search_input"
            placeholder="Search photos"
          />
        </div>

        <div className="header__profile header__item">
          {user.email && (
            <>
              <span className="header__storage_info">
                Storage:{' '}
                <code>
                  {bytesToMB(user.storage_used)}MB /{' '}
                  {bytesToGB(user.max_storage)}GB
                </code>
              </span>
              <div className="profile">
                <div className="profile__picture">
                  <AsyncImage src={user.picture} />
                </div>
                {user.username}
              </div>
              <ButtonSmall message="Logout" onClick={handleLogout} />
            </>
          )}
        </div>
      </header>
    </div>
  );
};

export default Header;
