import React from 'react';

import ImageIcon from '../../assets/images/icon-image.svg';

interface Props extends React.ComponentPropsWithRef<'div'> {
  onChange: React.ChangeEventHandler<HTMLInputElement> | undefined;
}

const UploadFileForm = (props: Props) => {
  const { onChange } = props;
  return (
    <div className="home__dropzone flex items-center justify-center w-full">
      <label
        htmlFor="dropzone-file"
        className="flex flex-col items-center justify-center w-full h-64 border-2 border-gray-300 border-dashed rounded-lg cursor-pointer bg-gray-50 dark:hover:bg-bray-800 dark:bg-gray-700 hover:bg-gray-100 dark:border-gray-600 dark:hover:border-gray-500 dark:hover:bg-gray-600"
      >
        <div className="flex flex-col items-center justify-center pt-5 pb-6">
          <img src={ImageIcon} alt="Upload Photo" className="h-24" />
          <p className="mb-2 text-sm text-gray-500 dark:text-gray-400">
            <span className="font-semibold">Click to upload</span> or drag and
            drop
          </p>
          <p className="text-xs text-gray-500 dark:text-gray-400">
            SVG, PNG, JPG or GIF
          </p>
        </div>
        <input
          id="dropzone-file"
          type="file"
          className="hidden"
          onChange={onChange}
          multiple
        />
      </label>
    </div>
  );
};

export default UploadFileForm;
