import 'bootstrap/dist/css/bootstrap.min.css';
import React, { useState, useEffect, useContext } from 'react';
import axios from 'axios';
import { UserContext } from './UserContext';

const Home = () => {
  const [file, setFile] = useState(null);
  const [search, setSearch] = useState('');
  const [files, setFiles] = useState([]);
  const { user } = useContext(UserContext);
  console.log(user, "user");
  const handleFileChange = (event) => {
    setFile(event.target.files[0]);
  };

  const handleUpload = async () => {
    if (!file) {
      alert('Please choose a file first.');
      return;
    }

    try {
      const response = await axios.post('http://localhost:8080/files/presignedurl', {
        key: user.id + '/' + file.name,
        action: 'upload',
      });

      console.log(response.data);
      const presignedUrl = response.data.presigned_url;
      console.log(presignedUrl);
      await axios.put(presignedUrl, file, {
        headers: {
          'Content-Type': file.type,
        },
      });
      alert('File uploaded successfully!');
      console.log('File uploaded successfully.');
      fetchFiles();
    } catch (error) {
      console.error('Error uploading file:', error);
    }
  };

  const handleSearchChange = (event) => {
    setSearch(event.target.value);
  };

  const handleSearch = async () => {
    try {
      const response = await axios.post('http://localhost:8080/files/search', {
        user_id: user.id,
        content: search,
      });
      setFiles(response.data);
    } catch (error) {
      console.error('Error searching files:', error);
    }
  };

  const handlePreview = async (file) => {
    try {
      console.log(file);
      const response = await axios.post('http://localhost:8080/files/redirecturl', {
        key : file.s3_key,
        action : 'read',
      });
      const url = response.data.presigned_url;
      window.open(url, '_blank'); 
      console.log(`Previewing file: ${file.file_name}`); 
    } catch (error) {
      console.error('Error previewing file:', error);
    }
  };

  const fetchFiles = async () => {
    if (!user || !user.id) {
      console.error('User is not defined');
      return;
    }

    try {
      const response = await axios.get(`http://localhost:8080/files/${user.id}`);
      setFiles(response.data);
    } catch (error) {
      console.error('Error fetching files:', error);
    }
  };

  useEffect(() => {
    if (user && user.id) {
      fetchFiles();
    }
  }, [user]);

  if (!user) {
    return <div>Loading...</div>;
  }

  return (
    <div className="container mt-5">
      <h2 className="text-primary mb-4">Home</h2>
      <div className="mb-3">
        <label htmlFor="fileUpload" className="form-label">
          Choose File
        </label>
        <input
          type="file"
          className="form-control"
          id="fileUpload"
          onChange={handleFileChange}
        />
        <button className="btn btn-primary mt-2" onClick={handleUpload}>
          Upload
        </button>
      </div>
      <div className="mb-3">
        <label htmlFor="searchField" className="form-label">
          Search Files
        </label>
        <input
          type="text"
          className="form-control"
          id="searchField"
          placeholder="Search"
          value={search}
          onChange={handleSearchChange}
        />
        <button className="btn btn-primary mt-2" onClick={handleSearch}>
          Search
        </button>
      </div>
      <div>
        <h4 className="mb-3">Files List</h4>
        <ul className="list-group">
          {files && files.map((file, index) => (
            <li key={index} className="list-group-item">
              <strong>Name:</strong> {file.file_name}<br />
              <strong>Size:</strong> {file.size} bytes<br />
              <strong>Type:</strong> {file.type}<br />
              <strong>Created At:</strong> {new Date(file.created_at).toLocaleString()}<br />
              <button className="btn btn-secondary mt-2" onClick={() => handlePreview(file)}>
                Download
              </button>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
};

export default Home;