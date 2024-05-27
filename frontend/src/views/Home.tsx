/* eslint-disable @typescript-eslint/no-unused-vars */
import React, { useEffect, useState } from 'react';
import fileStore from '../stores/file';
import { openDB, deleteDB } from 'idb';
import './Home.css';
import { MetaDTO } from '../types/file';

import AddIcon from '../assets/images/icon-add.svg';
import ButtonSmall from '../components/Buttons/ButtonSmall';

interface CacheDTO extends MetaDTO {
  imageUrl: string | null;
}

function blobToBase64(blob: Blob) {
  return new Promise((resolve, _) => {
    const reader = new FileReader();
    reader.onloadend = () => resolve(reader.result);
    reader.readAsDataURL(blob);
  });
}

const DB_NAME = 'fileCache';
const STORE_NAME = 'files';

const Home = () => {
  const { uploadFile, getMeta, getFile } = fileStore();

  const [images, setImages] = useState<CacheDTO[]>([]);
  const [loaded, setLoaded] = useState(false);
  const [uploadPopupOpen, setUploadPopupOpen] = useState(false);

  useEffect(() => {
    const fetchAndCacheImages = async () => {
      // Fetch metadata from server
      const metaList = await getMeta();
      if (!metaList) {
        setLoaded(true);
        return;
      }

      const db = await openDB(DB_NAME, 1, {
        upgrade(db) {
          db.createObjectStore(STORE_NAME, { keyPath: 'file_location' });
        },
      });

      // Get cached files from IndexedDB
      const cachedFiles = (await db.getAll(STORE_NAME)) as CacheDTO[];

      // Update cache
      const updatedCache = [...cachedFiles] as CacheDTO[];
      for (const meta of metaList) {
        if (
          !cachedFiles.some((file) => file.file_location === meta.file_location)
        ) {
          const file = { imageUrl: null, ...meta };
          updatedCache.push(file);
        }
      }

      // Fetch and cache new images
      for (const file of updatedCache) {
        if (!file.imageUrl) {
          const { blob } = await getFile(file.file_location);
          const base64 = await blobToBase64(blob);
          file.imageUrl = base64 as string;

          // Store the fetched image in IndexedDB
          await db.put(STORE_NAME, file);
        }
      }
      // Update state with cached images

      setImages(updatedCache);
      setLoaded(true);

      db.close();
    };
    fetchAndCacheImages();
  }, []);

  const groupedImages = images.reduce(
    (acc, file) => {
      const dateKey = (file.created_at as Date).toDateString();
      (acc[dateKey] ??= []).push(file); // Nullish coalescing operator for cleaner syntax
      return acc;
    },
    {} as Record<string, CacheDTO[]>,
  );

  const [files, setFiles] = useState<File[]>();
  const [uploadedStatus, setUploadedStatus] = useState(0);

  const handleFiles = async (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      const toUpload = [] as File[];
      for (let i = 0; i < e.target.files.length; i++) {
        const file = e.target.files[i];
        toUpload.push(file);
      }
      setFiles(toUpload);
    }
  };

  const handleUploadFiles = async () => {
    if (files) {
      setUploadedStatus(0);
      let uploaded = 0;
      const per = 100 / files.length;
      for (const file of files) {
        await uploadFile(file);
        uploaded += 1;
        setUploadedStatus(uploaded * per);
      }
      setUploadedStatus(0);
      setUploadPopupOpen(false);
    }
  };

  const handleOpenUploadPopup = () => {
    setUploadPopupOpen(true);
  };

  const handleCloseUploadPopup = () => {
    setUploadPopupOpen(false);
  };

  if (!loaded) {
    // Create skelleton loading for images
    return <div>Loading...</div>;
  }

  return (
    <div>
      <div className="home__images__wrapper">
        {Object.entries(groupedImages)
          .sort()
          .reverse()
          .map(([date, files]) => (
            <div key={date} className="home__images__block">
              <p>{date}</p>
              <div className="images__group">
                {files.map((file) => (
                  <img
                    key={file.file_location}
                    src={file.imageUrl || ''}
                    alt={file.file_name}
                    className="home__images"
                  />
                ))}
              </div>
            </div>
          ))}
      </div>

      <button className="home__open_upload" onClick={handleOpenUploadPopup}>
        <img src={AddIcon} alt="Upload Photos" />
        Upload Photos
      </button>
      <div
        className="home__upload__wrapper"
        style={{ display: uploadPopupOpen ? 'flex' : 'none' }}
      >
        <div className="home__upload">
          <p className="home__upload_label">Upload Photos</p>
          <div className="home__dropzone flex items-center justify-center w-full">
            <label
              htmlFor="dropzone-file"
              className="flex flex-col items-center justify-center w-full h-64 border-2 border-gray-300 border-dashed rounded-lg cursor-pointer bg-gray-50 dark:hover:bg-bray-800 dark:bg-gray-700 hover:bg-gray-100 dark:border-gray-600 dark:hover:border-gray-500 dark:hover:bg-gray-600"
            >
              <div className="flex flex-col items-center justify-center pt-5 pb-6">
                <svg
                  className="w-8 h-8 mb-4 text-gray-500 dark:text-gray-400"
                  aria-hidden="true"
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 20 16"
                >
                  <path
                    stroke="currentColor"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M13 13h3a3 3 0 0 0 0-6h-.025A5.56 5.56 0 0 0 16 6.5 5.5 5.5 0 0 0 5.207 5.021C5.137 5.017 5.071 5 5 5a4 4 0 0 0 0 8h2.167M10 15V6m0 0L8 8m2-2 2 2"
                  />
                </svg>
                <p className="mb-2 text-sm text-gray-500 dark:text-gray-400">
                  <span className="font-semibold">Click to upload</span> or drag
                  and drop
                </p>
                <p className="text-xs text-gray-500 dark:text-gray-400">
                  SVG, PNG, JPG or GIF
                </p>
              </div>
              <input
                id="dropzone-file"
                type="file"
                className="hidden"
                onChange={handleFiles}
                multiple
              />
            </label>
          </div>
          {files &&
            files.length <= 5 &&
            files?.map((f) => (
              <div key={f.name} className="home__preview">
                <img
                  src={URL.createObjectURL(f)}
                  className="home__image_preview"
                />
                <p>{f.name}</p>
              </div>
            ))}
          {files && files.length > 5 && (
            <p className="home__preview_selected">
              Selected {files.length} photos
            </p>
          )}

          {uploadedStatus == 0 ? (
            <div className="home__popup_buttons">
              <ButtonSmall
                message="Close"
                onClick={handleCloseUploadPopup}
                className="home__close_popup"
              ></ButtonSmall>
              <ButtonSmall
                message="Upload Images"
                onClick={handleUploadFiles}
              ></ButtonSmall>
            </div>
          ) : (
            <p style={{ textAlign: 'center' }}>
              Loading... {uploadedStatus.toFixed(1)}%
            </p>
          )}
        </div>
      </div>
    </div>
  );
};

export default Home;
