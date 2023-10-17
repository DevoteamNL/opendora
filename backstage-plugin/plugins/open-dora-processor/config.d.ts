export interface Config {
  'open-dora': {
    /**
     * OpenDORA DevLake root URL
     * @visibility backend
     */
    devLakeBaseUrl: string;
    /**
     * OpenDORA DevLake GitLab connection ID
     * @visibility backend
     */
    devLakeConnectionId: number;
    /**
     * OpenDORA DevLake username/password authorization
     * @visibility secret
     */
    devLakeBasicAuth: string;
  };
}
