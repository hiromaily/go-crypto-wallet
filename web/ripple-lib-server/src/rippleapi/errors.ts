
export interface rippledError {
  data: errorContent;
}

export interface errorContent {
  error: string;
  error_code: number;
  error_message: string;
  request: object;
  status: string;
  type: string;
  validated: boolean;
}
