// 
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
// 
//       http://www.apache.org/licenses/LICENSE-2.0
// 
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

// Bassa js client library will enable you to interact easily with Bassa API server with most of its functions covered. 

const axios = require('axios');
const { URL } = require('url');
const { InvalidUrl, IncompleteParams, ResponseError } = require('./errors');

class Bassa {
    constructor(apiUrl, total = 1, backoffFactor = 1, timeout = 5) {
        this.apiUrl = apiUrl;
        this.headers = {
            'Content-Type': 'application/x-www-form-urlencoded',
        };

        const reg = new RegExp(
            '^(?:http|ftp)s?://' +  // http:// or https://
            // domain...
            '(?:(?:[A-Z0-9](?:[A-Z0-9-]{0,61}[A-Z0-9])?\\.)+(?:[A-Z]{2,6}\\.?|[A-Z0-9-]{2,}\\.?)|' +
            'localhost|' +  // localhost...
            '\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3})' +  // ...or ip
            '(?::\\d+)?' +  // optional port
            '(?:/?|[/?]\\S+)$', 'i'
        );

        if (!reg.test(apiUrl)) {
            throw new InvalidUrl('Invalid URL');
        }

        const retryStrategy = {
            total,
            backoffFactor,
            retry: [429, 500, 502, 503, 504],
        };

        const instance = axios.create({
            baseURL: apiUrl,
            timeout: timeout * 1000,
            httpAgent: new axios.Agent({
                http: new axios.Agent(retryStrategy),
                https: new axios.Agent(retryStrategy),
            }),
        });

        this.http = instance;
    }

    // User Functions

    login(user_name, password) {
        const endpoint = "/api/login";
        const url = new URL(endpoint, this.apiUrl);
        const params = new URLSearchParams();
        if (!user_name || !password) {
            throw new IncompleteParams();
        }
        params.append('user_name', user_name);
        params.append('password', password);

        return this.http.post(url.href, params, { headers: this.headers })
            .then(response => {
                if (response.status === 200) {
                    this.headers['token'] = response.headers.token;
                } else {
                    throw new ResponseError(`API response: ${response.status}`);
                }
            });
    }

    addRegularUserRequest(user_name, password, email) {
        const endpoint = "/api/regularuser";
        const url = new URL(endpoint, this.apiUrl);
        const params = new URLSearchParams();
        if (!user_name || !password || !email) {
            throw new IncompleteParams();
        }
        params.append('user_name', user_name);
        params.append('password', password);
        params.append('email', email);
    
        return this.http.post(url.href, params, { headers: this.headers });
    }

    addUserRequest(user_name, password, email, auth_level = 1) {
        const endpoint = "/api/user";
        const url = new URL(endpoint, this.apiUrl);
        const params = new URLSearchParams();
        if (!user_name || !password || !email) {
            throw new IncompleteParams();
        }
        params.append('user_name', user_name);
        params.append('password', password);
        params.append('email', email);
        params.append('auth', auth_level);
    
        return this.http.post(url.href, params, { headers: this.headers });
    }

    removeUserRequest(user_name) {
        if (!user_name) {
            throw new IncompleteParams();
        }
        const endpoint = "/api/user";
        const url = new URL(`${endpoint}/${user_name}`, this.apiUrl);
    
        return this.http.delete(url.href, { headers: this.headers });
    }

    updateUserRequest(user_name, new_user_name, password, auth_level, email) {
        if (!user_name || !new_user_name || !password || auth_level == null || !email) {
            throw new IncompleteParams();
        }
        const endpoint = "/api/user";
        const url = new URL(`${endpoint}/${user_name}`, this.apiUrl);
        const params = new URLSearchParams();
        params.append('user_name', new_user_name);
        params.append('password', password);
        params.append('auth_level', auth_level);
        params.append('email', email);
    
        return this.http.put(url.href, params, { headers: this.headers });
    }

    getUserRequest() {
        const endpoint = "/api/user";
        const url = new URL(endpoint, this.apiUrl);
    
        return this.http.get(url.href, { headers: this.headers })
            .then(response => {
                if (response.status === 200) {
                    return response.json();
                } else {
                    throw new Error(`Request failed with status ${response.status}`);
                }
            });
    }

    getUserSignupRequests() {
        const endpoint = "/api/user/requests";
        const url = new URL(endpoint, this.apiUrl);
    
        return this.http.get(url.href, { headers: this.headers })
            .then(response => {
                if (response.status === 200) {
                    return response.json();
                } else {
                    throw new Error(`Request failed with status ${response.status}`);
                }
            });
    }

    approveUserRequest(user_name) {
        if (!user_name) {
            throw new IncompleteParams();
        }
        const endpoint = "/api/user/approve";
        const url = new URL(`${endpoint}/${user_name}`, this.apiUrl);
    
        return this.http.post(url.href, null, { headers: this.headers });
    }

    getBlockedUsersRequest() {
        const endpoint = "/api/user/blocked";
        const url = new URL(endpoint, this.apiUrl);
    
        return this.http.get(url.href, { headers: this.headers })
            .then(response => {
                if (response.status === 200) {
                    return response.json();
                } else {
                    throw new Error(`Request failed with status ${response.status}`);
                }
            });
    }

    blockUserRequest(user_name) {
        if (!user_name) {
            throw new IncompleteParams();
        }
        const endpoint = "/api/user/blocked";
        const url = new URL(`${endpoint}/${user_name}`, this.apiUrl);
    
        return this.http.post(url.href, null, { headers: this.headers });
    }

    unblockUserRequest(user_name) {
        if (!user_name) {
            throw new IncompleteParams();
        }
        const endpoint = "/api/user/blocked";
        const url = new URL(`${endpoint}/${user_name}`, this.apiUrl);
    
        return this.http.delete(url.href, { headers: this.headers });
    }

    getDownloadsUserRequest(limit = 1) {
        const endpoint = "/api/user/downloads";
        const url = new URL(`${endpoint}/${limit}`, this.apiUrl);
    
        return this.http.get(url.href, { headers: this.headers })
            .then(response => {
                if (response.status === 200) {
                    return response.json();
                } else {
                    throw new Error(`Request failed with status ${response.status}`);
                }
            });
    }

    getToptenHeaviestUsers() {
        const endpoint = "/api/user/heavy";
        const url = new URL(endpoint, this.apiUrl);
    
        return this.http.get(url.href, { headers: this.headers })
            .then(response => {
                if (response.status === 200) {
                    return response.json();
                } else {
                    throw new Error(`Request failed with status ${response.status}`);
                }
            });
    }

    // Download Functions

    startDownload(server_key = "123456789") {
        const endpoint = "/api/download/start";
        const url = new URL(endpoint, this.apiUrl);
    
        this.headers['key'] = server_key;
    
        return this.http.get(url.href, { headers: this.headers })
            .then(response => {
                if (response.status === 200) {
                    return response.json();
                } else {
                    throw new Error(`Request failed with status ${response.status}`);
                }
            });
    }

    killDownload(server_key = "123456789") {
        const endpoint = "/api/download/kill";
        const url = new URL(endpoint, this.apiUrl);
    
        this.headers['key'] = server_key;
    
        return this.http.get(url.href, { headers: this.headers })
            .then(response => {
                if (response.status === 200) {
                    return response.json();
                } else {
                    throw new Error(`Request failed with status ${response.status}`);
                }
            });
    }

    addDownloadRequest(download_link) {
        if (!download_link) {
            throw new IncompleteParams();
        }
        const endpoint = "/api/download";
        const url = new URL(endpoint, this.apiUrl);
        const params = { link: download_link };
    
        return this.http.post(url.href, JSON.stringify(params), { headers: this.headers })
            .then(response => {
                if (response.status !== 200) {
                    throw new Error("Add download was not successful");
                }
                return response.json();
            });
    }

    removeDownloadRequest(id) {
        if (!id) {
            throw new Error("ID is required to remove download");
        }
        const endpoint = `/api/download/${id}`;
        const url = new URL(endpoint, this.apiUrl);
    
        return this.http.delete(url.href, { headers: this.headers })
            .then(response => {
                if (response.status !== 200) {
                    throw new Error(`Remove download failed with status ${response.status}`);
                }
                return response.json();
            });
    }

    rateDownloadRequest(id, rate) {
        if (!id || !rate) {
            throw new Error("ID and rate are required to rate download");
        }
        const endpoint = `/api/download/${id}`;
        const url = new URL(endpoint, this.apiUrl);
        const params = { rate: rate };
    
        return this.http.post(url.href, JSON.stringify(params), { headers: this.headers })
            .then(response => {
                if (response.status !== 200) {
                    throw new Error(`Rate download failed with status ${response.status}`);
                }
                return response.json();
            });
    }

    getDownloadsRequest(limit) {
        if (limit === undefined || limit === null) {
            throw new Error("Limit is required to get downloads");
        }
        const endpoint = `/api/downloads/${limit}`;
        const url = new URL(endpoint, this.apiUrl);
    
        return this.http.get(url.href, { headers: this.headers })
            .then(response => {
                if (response.status === 200) {
                    return response.json();
                } else {
                    throw new Error(`Get downloads failed with status ${response.status}`);
                }
            });
    }

    getDownload(id) {
        if (!id) {
            throw new Error("ID is required to get download");
        }
        const endpoint = `/api/download/${id}`;
        const url = new URL(endpoint, this.apiUrl);
    
        return this.http.get(url.href, { headers: this.headers })
            .then(response => {
                if (response.status === 200) {
                    return response.json();
                } else {
                    throw new Error(`Get download failed with status ${response.status}`);
                }
            });
    }

    // File Functions

    startCompression(gidList) {
        if (!gidList || gidList.length === 0) {
            throw new Error("GID list is required to start compression");
        }
        const endpoint = "/api/compress";
        const url = new URL(endpoint, this.apiUrl);
        const params = { gid: gidList };
    
        return this.http.post(url.href, JSON.stringify(params), { headers: this.headers })
            .then(response => {
                if (response.status === 200) {
                    return response.json();
                } else {
                    throw new Error(`Start compression failed with status ${response.status}`);
                }
            });
    }

    getCompressionProgress(id) {
        if (!id) {
            throw new Error("ID is required to get compression progress");
        }
        const endpoint = `/api/compression-progress/${id}`;
        const url = new URL(endpoint, this.apiUrl);
    
        return this.http.get(url.href, { headers: this.headers })
            .then(response => {
                if (response.status === 200) {
                    return response.json();
                } else {
                    throw new Error(`Get compression progress failed with status ${response.status}`);
                }
            });
    }

    sendFileFromPath(id) {
        if (!id) {
            throw new Error("ID is required to send file from path");
        }
        const endpoint = "/api/file";
        const url = new URL(endpoint, this.apiUrl);
        url.searchParams.append('gid', id);
    
        return this.http.get(url.href, { headers: this.headers })
            .then(response => {
                if (response.status === 200) {
                    return response.json();
                } else {
                    throw new Error(`Send file from path failed with status ${response.status}`);
                }
            });
    }
    
}

module.exports = Bassa