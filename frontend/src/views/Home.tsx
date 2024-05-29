/* eslint-disable indent */
/* eslint-disable @typescript-eslint/no-unused-vars */
import React, { useEffect, useState } from 'react';
import fileStore from '../stores/file';
import { openDB, deleteDB } from 'idb';
import './Home.css';
import { MetaDTO } from '../types/file';

import AddIcon from '../assets/images/icon-add.svg';
import ButtonSmall from '../components/Buttons/ButtonSmall';
import DeleteIcon from '../assets/images/icon-delete.svg';
import UploadFileForm from '../components/Forms/UploadFile';
import classNames from 'classnames';

interface CacheDTO extends MetaDTO {
  imageUrl: string | null;
  isChecked: boolean;
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
  const { uploadFile, getMeta, getFile, deleteFile } = fileStore();

  const [refreshTrigger, setRefreshTrigger] = useState(0);
  const [images, setImages] = useState<CacheDTO[]>([]);
  const [loaded, setLoaded] = useState(false);
  const [uploadPopupOpen, setUploadPopupOpen] = useState(false);

  useEffect(() => {
    const fetchAndCacheImages = async () => {
      // Fetch metadata from the server
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

      // Create sets for efficient comparison
      const metaFileLocations = new Set(
        metaList.map((meta) => meta.file_location),
      );
      const cachedFileLocations = new Set(
        cachedFiles.map((file) => file.file_location),
      );

      // Identify new files to add and missing files to remove
      const filesToAdd = metaList.filter(
        (meta) => !cachedFileLocations.has(meta.file_location),
      );
      const filesToRemove = cachedFiles.filter(
        (file) => !metaFileLocations.has(file.file_location),
      );

      // Update cache and IndexedDB
      for (const file of filesToRemove) {
        await db.delete(STORE_NAME, file.file_location);
        const indexToRemove = cachedFiles.indexOf(file); // Find index for removal
        if (indexToRemove > -1) {
          cachedFiles.splice(indexToRemove, 1); // Remove the item
        }
      }

      for (const meta of filesToAdd) {
        const file = { imageUrl: '', ...meta, isChecked: false } as CacheDTO;
        cachedFiles.push(file); // Add to cache before fetching

        try {
          const { blob } = await getFile(meta.file_location);
          const base64 = await blobToBase64(blob);
          file.imageUrl = base64 as string;
          await db.put(STORE_NAME, file); // Store in IndexedDB
        } catch (error) {
          console.error('Error fetching or storing image:', error);
          // Consider removing from cache if fetch fails
          cachedFiles.pop(); // Remove the last added file (the one that failed)
        }
      }

      // Update state with cached images
      setImages(cachedFiles); // Update state after modifications
      setLoaded(true);

      db.close();
    };

    fetchAndCacheImages();
  }, [refreshTrigger]);

  // Group images by date
  const [groupedImages, setGroupedImages] = useState<
    Record<string, CacheDTO[]>
  >({});
  useEffect(() => {
    console.log('images updated');
    const grouped = images.reduce(
      (acc, image) => {
        const dateKey = (image.created_at as Date).toDateString();
        (acc[dateKey] ??= []).push(image); // Nullish coalescing operator for cleaner syntax
        return acc;
      },
      {} as Record<string, CacheDTO[]>,
    );
    setGroupedImages(grouped);
  }, [images]);

  const [groupImageCheck, setGroupImagesCheck] = useState<{
    [key: string]: string;
  }>({});
  useEffect(() => {
    const initialGroupImageCheck = Object.keys(groupedImages).reduce(
      (acc: { [key: string]: string }, key) => {
        const files = groupedImages[key];
        const allChecked = files.every((file) => file.isChecked);
        const someChecked = files.some((file) => file.isChecked);
        acc[key] = allChecked
          ? 'checked'
          : someChecked
            ? 'half-checked'
            : 'unchecked';
        return acc;
      },
      {},
    );
    setGroupImagesCheck(initialGroupImageCheck);
  }, [groupedImages]);

  const handleImageCheck = (location: string) => {
    setImages(() =>
      images.map((image) => {
        if (image.file_location === location) {
          return { ...image, isChecked: !image.isChecked };
        }
        return { ...image };
      }),
    );
  };

  const handleImageGroupCheck = (date: string) => {
    const filtered = images.filter(
      ({ created_at }) => (created_at as Date).toDateString() === date,
    );

    const allChecked = filtered.every((image) => image.isChecked);
    setImages(() =>
      images.map((image) => {
        const createdAt = image.created_at as Date;
        if (createdAt.toDateString() === date) {
          return {
            ...image,
            isChecked: allChecked ? false : true,
          };
        }
        return { ...image };
      }),
    );
  };

  const handleCancelImageSelect = () => {
    setImages(() =>
      images.map((image) => {
        return {
          ...image,
          isChecked: false,
        };
      }),
    );
  };

  const handleDeleteSelectedImages = () => {
    const toDelete = images.filter((image) => image.isChecked);
    for (const image of toDelete) {
      handleDeleteFile(image.id);
    }
  };

  // Uploading files manipulations
  const [files, setFiles] = useState<File[]>();
  const [uploadedStatus, setUploadedStatus] = useState(0);

  const handleSelectFiles = async (e: React.ChangeEvent<HTMLInputElement>) => {
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
      setRefreshTrigger((prev) => prev + 1);
    }
  };

  const handleOpenUploadPopup = () => {
    setUploadPopupOpen(true);
  };

  const handleCloseUploadPopup = () => {
    setUploadPopupOpen(false);
  };

  const handleDeleteFile = async (id: string) => {
    await deleteFile(id);
    setRefreshTrigger((prev) => prev + 1);
  };

  if (!loaded) {
    // Create skelleton loading for images
    return <div>Loading...</div>;
  }

  return (
    <div>
      {images.some((image) => image.isChecked) ? (
        <div className="home__selected home__selected--active">
          <div className="home__selected_left">
            <ButtonSmall message="cancel" onClick={handleCancelImageSelect} />
            <p>{images.filter((image) => image.isChecked).length} selected</p>
          </div>
          <div>
            <ButtonSmall
              message="delete"
              onClick={handleDeleteSelectedImages}
            />
          </div>
        </div>
      ) : (
        <div className="home__selected"></div>
      )}

      <div className="home__images__wrapper">
        {Object.entries(groupedImages)
          .sort()
          .reverse()
          .map(([date, files]) => {
            let groupCheckStyle: string;
            switch (groupImageCheck[date]) {
              case 'checked':
                groupCheckStyle = 'home__group_input--checked';
                break;
              case 'half-checked':
                groupCheckStyle = 'home__group_input--half-checked';
                break;
              default:
                groupCheckStyle = 'home__group_input--unchecked';
                break;
            }

            return (
              <div key={date} className="home__images_group">
                <div>
                  <label className="home__group_label">
                    {date}
                    <input
                      type="checkbox"
                      name="checked"
                      className={classNames(
                        groupCheckStyle,
                        'home__group_input',
                      )}
                      onChange={() => handleImageGroupCheck(date)}
                    />
                    <span className="home__group_checkbox"></span>
                  </label>
                </div>
                <div className="images__group">
                  {files.map((file) => (
                    <div key={file.file_location} className="home__images">
                      <img src={file.imageUrl || ''} alt={file.file_name} />
                      <div className="home__image_hover">
                        <label className="home__image_label">
                          <input
                            type="checkbox"
                            name="checked"
                            checked={file.isChecked}
                            className="home__image_input"
                            onChange={() =>
                              handleImageCheck(file.file_location)
                            }
                          />
                          <span className="home__image_checkbox"></span>
                        </label>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            );
          })}
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
          <UploadFileForm onChange={handleSelectFiles} />
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
