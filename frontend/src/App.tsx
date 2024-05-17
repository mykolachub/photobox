import React from 'react';
import './App.css';
import { Route, Routes } from 'react-router-dom';
import Layout from './components/Layout';
import Signup from './views/Signup';
import Login from './views/Login';
import ProtectedRoute from './components/ProtectedRoute';
import Home from './views/Home';

function App() {
  return (
    <Routes>
      <Route path="/" element={<Layout />}>
        {/* public routes */}
        <Route path="signup" element={<Signup />} />
        <Route path="callback" element={<Login />} />

        {/* private routes */}
        <Route path="" element={<ProtectedRoute />}>
          <Route path="/" element={<Home />} />
        </Route>
      </Route>
    </Routes>
  );
}

export default App;
