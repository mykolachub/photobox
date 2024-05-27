export interface MetaDTO {
  id: string;
  user_id: string;
  file_location: string;
  file_name: string;
  file_ext: string;
  file_size: number;
  file_last_modified: Date | ProtoTime;
  created_at: Date | ProtoTime;
}

export interface FileDTO {
  file_location: string;
  url: string;
}

export interface ProtoTime {
  seconds: number;
  nanos: number;
}
