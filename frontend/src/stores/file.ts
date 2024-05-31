import { create } from 'zustand';
import { FileDTO, MetaDTO, ProtoTime } from '../types/file';
import { GetMetaResponse, ServerResponse } from '../types/server';
import axios from 'axios';
import config from '../utils/config';
import authStore from './auth';

const API_URL = config.env.apiUrl;

interface FileState {
  meta: MetaDTO[];
  files: FileDTO[];
  uploadFile: (file: File) => Promise<MetaDTO>;
  getMeta: () => Promise<MetaDTO[]>;
  getFile: (file_location: string) => Promise<Test>;
  deleteFile: (id: string) => Promise<void>;
}

interface Test {
  url: string;
  blob: Blob;
}

const fileStore = create<FileState>((set) => ({
  meta: [],
  files: [],
  async uploadFile(file: File): Promise<MetaDTO> {
    try {
      const formData = new FormData();
      formData.append('file', file);
      formData.append('lastModified', file.lastModified.toString());

      const reader = new FileReader();
      const { width, height } = await new Promise<{
        width: number;
        height: number;
      }>((resolve) => {
        reader.onload = () => {
          const imageDataUrl = reader.result as string;
          const image = new Image();
          image.src = imageDataUrl;

          image.onload = () => {
            const width = image.width;
            const height = image.height;
            resolve({ width, height });
          };
        };

        reader.readAsDataURL(file);
      });

      formData.append('fileWidth', width.toString());
      formData.append('fileHeight', height.toString());

      const token = `Bearer ${authStore.getState().token}`;
      const response = await axios.post(API_URL + '/meta', formData, {
        headers: { Authorization: token },
      });
      const { data } = response.data as ServerResponse<MetaDTO>;

      return data;
    } catch (error) {
      throw new Error('' + error);
    }
  },
  async getMeta(): Promise<MetaDTO[]> {
    try {
      const token = `Bearer ${authStore.getState().token}`;
      const response = await axios.get(API_URL + '/meta', {
        headers: { Authorization: token },
      });
      const { data } = response.data as ServerResponse<GetMetaResponse>;

      // Convert Protobuf date to Date
      if (!data.metas) {
        return data.metas;
      }

      data.metas = data.metas.map(
        ({ file_last_modified, created_at, ...rest }) => ({
          file_last_modified: protoToDate(file_last_modified as ProtoTime),
          created_at: protoToDate(created_at as ProtoTime),
          ...rest,
        }),
      );

      set({ meta: [...data.metas] });
      return data.metas;
    } catch (error) {
      throw new Error('' + error);
    }
  },
  async getFile(file_location: string): Promise<Test> {
    try {
      const token = `Bearer ${authStore.getState().token}`;
      const response = await axios.get(
        API_URL + '/meta/files?file_location=' + file_location,
        { headers: { Authorization: token }, responseType: 'blob' },
      );
      const url = URL.createObjectURL(response.data);
      set((state) => ({
        files: [...state.files, { file_location, url }],
      }));
      return { url, blob: response.data };
    } catch (error) {
      throw new Error('' + error);
    }
  },
  async deleteFile(id: string): Promise<void> {
    try {
      const token = `Bearer ${authStore.getState().token}`;
      const response = await axios.delete(API_URL + `/meta/${id}`, {
        headers: { Authorization: token },
      });
      const { data } = response.data as ServerResponse<void>;
      console.log(data);
      return data;
    } catch (error) {
      throw new Error('' + error);
    }
  },
}));

export default fileStore;

const protoToDate = (time: ProtoTime): Date => {
  const milliseconds = time.seconds * 1000 + time.nanos / 1000000;
  return new Date(milliseconds);
};
