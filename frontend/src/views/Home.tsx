import axios from 'axios';
import React from 'react';
import authStore from '../stores/auth';

const Home = () => {
  const { token } = authStore();
  const handleFile = async (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      const file = e.target.files[0];
      console.log(file);

      const formData = new FormData();
      formData.append('file', file);
      formData.append('lastModified', file.lastModified.toString());

      try {
        const sent = await axios.post(
          'http://localhost:8080/api/meta',
          formData,
          {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          },
        );
        console.log(sent);
      } catch (error) {
        console.log(error);
      }
    }
  };
  return (
    <div>
      <input type="file" name="file" id="file" onChange={handleFile} />
    </div>
  );
};

export default Home;
