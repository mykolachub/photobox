/* eslint-disable @typescript-eslint/no-unused-vars */
/* eslint-disable indent */
import React, { useEffect, useState } from 'react';
import fileStore from '../stores/file';
import './Home.css';
import { MetaDTO } from '../types/file';
import { useInView } from 'react-intersection-observer';

import AddIcon from '../assets/images/icon-add.svg';
import ButtonSmall from '../components/Buttons/ButtonSmall';
import ImageIcon from '../assets/images/icon-image.svg';
import UploadFileForm from '../components/Forms/UploadFile';
import classNames from 'classnames';
import searchStore from '../stores/search';

interface CacheDTO extends MetaDTO {
  imageUrl: string | null;
  isChecked: boolean;
}

function blobToBase64(blob: Blob) {
  return new Promise((resolve) => {
    const reader = new FileReader();
    reader.onloadend = () => resolve(reader.result);
    reader.readAsDataURL(blob);
  });
}

interface ImageComponentProps {
  file: CacheDTO;
  image: string;
  setImages: React.Dispatch<React.SetStateAction<CacheDTO[]>>;
}

const ImageComponent: React.FC<ImageComponentProps> = ({
  file,
  image,
  setImages,
}) => {
  const { getFile } = fileStore();
  const [imageUrl, setImageUrl] = useState<string | null>(file.imageUrl);
  const [loading, setLoading] = useState(false);

  const { ref, inView } = useInView({ threshold: 0 });

  useEffect(() => {
    if (file.imageUrl === '' && inView) {
      const fetchImage = async () => {
        setLoading(true);
        try {
          const { url } = await getFile(file.file_location);
          setImageUrl(url);
          setImages((prev) =>
            prev.map((prevImage) => {
              if (prevImage.id === file.id) {
                return { ...prevImage, imageUrl: url };
              }
              return { ...prevImage };
            }),
          );
        } catch (error) {
          console.error(error);
        } finally {
          setLoading(false);
        }
      };
      fetchImage();
    }
  }, [ref, inView]);

  const width = (200 * file.file_width) / file.file_height;
  const height = 200;

  return (
    <div
      ref={ref}
      style={{
        width: width,
        height: height,
      }}
    >
      {loading ? (
        <div
          style={{
            width: '100%',
            height: '100%',
            background: '#f2f2f2',
          }}
        ></div>
      ) : imageUrl ? (
        <img src={imageUrl} alt={file.file_name} />
      ) : (
        <div
          style={{
            width: '100%',
            height: '100%',
            background: '#f2f2f2',
          }}
        ></div>
      )}
    </div>
  );
};
const Home = () => {
  const { uploadFile, getMeta, deleteFile } = fileStore();
  const { search } = searchStore();

  const [refreshMetaTrigger, setRefreshMetaTrigger] = useState(0);
  const [refreshSearchInputTrigger, setRefreshSearchInputTrigger] = useState(0);

  const [images, setImages] = useState<CacheDTO[]>([]);
  const [imagesCopy, setImagesCopy] = useState<CacheDTO[]>([]);
  const [indexLabelToImage, setIndexLabelToImage] = useState<{
    [key: string]: CacheDTO[];
  }>({});

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

      const cachedFiles = [] as CacheDTO[];

      for (const meta of metaList) {
        const file = { imageUrl: '', ...meta, isChecked: false } as CacheDTO;
        cachedFiles.push(file);
      }

      // Update state with cached images
      setImages(cachedFiles); // Update state after modifications
      setLoaded(true);
    };

    fetchAndCacheImages();
  }, [refreshMetaTrigger]);

  // when images changed, imagesCopy will be affected to
  // imagesCopy is for image grouping, searching
  useEffect(() => {
    console.log('images changed');
    // Creating Indexing object Label to Images
    const newIndexLabelToImage: { [key: string]: CacheDTO[] } = {};
    images.forEach((image) => {
      image.labels.forEach((label) => {
        const value = label.value.toLowerCase();
        if (!newIndexLabelToImage[value]) {
          newIndexLabelToImage[value] = [];
        }
        newIndexLabelToImage[value].push(image);
      });
    });
    setIndexLabelToImage(newIndexLabelToImage);
    if (search) {
      setRefreshSearchInputTrigger((prev) => prev + 1);
    }
    setImagesCopy(images);
  }, [images]);

  // Searching images by label
  useEffect(() => {
    if (!search) {
      setImagesCopy(images);
      return;
    }
    const matchingMetas: CacheDTO[] = [];

    Object.keys(indexLabelToImage).forEach((label) => {
      if (
        label.startsWith(search) || // starts with
        label.endsWith(search) || // ends with
        label.includes(search) // contains
      ) {
        matchingMetas.push(...indexLabelToImage[label]);
      }
    });

    const uniqueMatchingMetas = [...new Set(matchingMetas)];
    setImagesCopy(uniqueMatchingMetas);
  }, [search, refreshSearchInputTrigger]);

  // Group images by date
  const [groupedImages, setGroupedImages] = useState<
    Record<string, CacheDTO[]>
  >({});
  useEffect(() => {
    const grouped = imagesCopy.reduce(
      (acc, image) => {
        const dateKey = (image.created_at as Date).toDateString();
        (acc[dateKey] ??= []).push(image); // Nullish coalescing operator for cleaner syntax
        return acc;
      },
      {} as Record<string, CacheDTO[]>,
    );
    setGroupedImages(grouped);
  }, [imagesCopy]);

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
    console.log(images.find(({ file_location }) => file_location === location));
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
    if (search) {
      const t = groupedImages[date].map((image) => image.id);
      const allChecked = groupedImages[date].every((image) => image.isChecked);
      setImages(() =>
        images.map((image) => {
          if (t.includes(image.id)) {
            return {
              ...image,
              // isChecked: true,
              isChecked: allChecked ? false : true,
            };
          }
          return { ...image };
        }),
      );
      return;
    }

    const filtered = imagesCopy.filter(
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
      setRefreshMetaTrigger((prev) => prev + 1);
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
    setRefreshMetaTrigger((prev) => prev + 1);
  };

  if (!loaded) {
    // Create skelleton loading for images
    return <div>Loading...</div>;
  }

  return (
    <div className="home__content">
      <aside className="home__left">
        <div className="home__datelinks">
          {Object.keys(groupedImages)
            .sort(
              (a, b) =>
                new Date(a).getMilliseconds() - new Date(b).getMilliseconds(),
            )
            .reverse()
            .map((date) => (
              <a
                key={date}
                href={'#' + date.replaceAll(' ', '_')}
                className="home__date_link"
              >
                {'- ' + date}
              </a>
            ))}
        </div>
      </aside>

      <div className="home__main">
        {images.some((image) => image.isChecked) ? (
          <div className="home__selected_wrapper">
            <div className="home__selected home__selected--active">
              <div className="home__selected_left">
                <ButtonSmall
                  message="cancel"
                  onClick={handleCancelImageSelect}
                />
                <p>
                  {images.filter((image) => image.isChecked).length} selected
                </p>
              </div>
              <div className="home__selected_images">
                {images.map((image) =>
                  image.isChecked ? (
                    <img
                      key={image.id}
                      src={image.imageUrl as string}
                      className="home__selected_image"
                    />
                  ) : (
                    <></>
                  ),
                )}
              </div>
              <div>
                <ButtonSmall
                  message="delete"
                  onClick={handleDeleteSelectedImages}
                />
              </div>
            </div>
          </div>
        ) : (
          <div className="home__selected"></div>
        )}

        <div className="home__images__wrapper">
          {Object.keys(groupedImages).length ? (
            Object.keys(groupedImages)
              .sort(
                (a, b) =>
                  new Date(a).getMilliseconds() - new Date(b).getMilliseconds(),
              )
              .reverse()
              .map((date) => {
                const files = groupedImages[date];
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
                    <div id={date.replaceAll(' ', '_')}>
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
                      {files.map((file) => {
                        return (
                          <div
                            key={file.file_location}
                            className="home__images"
                          >
                            <ImageComponent
                              file={file}
                              image={file.imageUrl ?? ''}
                              setImages={setImages}
                            />
                            {/* <img /> */}
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
                        );
                      })}
                    </div>
                  </div>
                );
              })
          ) : (
            <div className="home__noimages">
              <img
                src={ImageIcon}
                alt="No Photos"
                className="home__noimages_icon"
              />
              {search ? (
                <p>No images found</p>
              ) : (
                <>
                  <p>No images yet</p>
                  <ButtonSmall
                    message="Upload Photos"
                    onClick={handleOpenUploadPopup}
                  />
                </>
              )}
            </div>
          )}
        </div>
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
