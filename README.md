# 🗂️ S3-Compatible File Upload and Download Service (Go + AWS SDK v2)

This Go application provides a complete S3-compatible client to upload and download files using AWS S3. It enforces strict policies on bucket creation and object uploads, mimicking production-grade file storage behavior.



## 📦 Features

✅ Automatically creates a **daily bucket** named `<team>-YYYYMMDD`  
✅ Enforces three core policies:
1. Bucket name must start with team name
2. Versioning must be enabled
3. Total bucket size (including all versions) must not exceed **1GB**

✅ Uploads files:
- Only if there are at least **5 files**
- Each file must be **≥10MB**
- Uploads are run **in parallel**

✅ Downloads multiple files to a `downloads/` folder  
✅ Uses **AWS SDK v2** and supports both **AWS S3** and **MinIO**



## ⚙️ Local Setup

This Go client connects to a **MinIO server running locally**, accessible at:  
**http://localhost:9001**

You must first run MinIO using Docker Compose or any container engine exposing port `9000` for API and `9001` for console UI.



## 📁 Folder Structure
```
├── cmd
│   └── main.go
├── downloads
├── files
├── go.mod
├── go.sum
├── pkg
│   ├── client
│   │   ├── client.go
│   │   └── type.go
│   └── storage
│       ├── bucket.go
│       ├── file.go
│       └── utils.go
└── README.md
```


## 🚀 How to Run the Application

### 1. Clone the Repository

```bash
git clone git@github.com:georgeroumieh-pon/s3-client.git
cd s3-client
```

### 2. Set Required Environment Variables

Before running, export the required MinIO credentials:

```bash
export ACCESS_KEY=<access_key>
export SECRET_KEY=<secret_key>
```

### 3. Create Folders

In the root of the project, create:

```bash
mkdir files downloads
```

- Place at least 5 files (≥10MB) in the `files/` folder
- Downloaded files will be saved into the `downloads/` folder

### 4. Run the App

```bash
go run ./cmd
```

This will:
- Create a new daily bucket
- Enable versioning
- Upload all files from the `files/` folder
- Download a list of predefined files into `downloads/`
