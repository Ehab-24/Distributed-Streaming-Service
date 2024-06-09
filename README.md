## Running Locally

### 1. Clone the project
```
git clone https://github.com/Ehab-24/Distributed-Streaming-Service.git
cd Distributed-Streaming-Service
```

### 2. Install Bento4 tools
To install the [Bento4 tools](https://github.com/axiomatic-systems/Bento4), execute the following commands:
```
git clone https://github.com/axiomatic-systems/Bento4.git
cd Bento4/

mkdir cmakebuild
cd cmakebuild/
cmake -DCMAKE_BUILD_TYPE=Release ..
make

sudo make install
```

### 3. Initialize the master
```
cd apps/master
python -m venv venv
source venv/bin/activate
pip install -r requirements.txt
```

Apply database migrations:
```
python manage.py makemigrations
python manage.py migrate
```

Create a super user:
```
python manage.py createsuperuser
```

Now, to add chunk server details, start the master:
```
python manage.py runserver
```
open the [django admin page](http://localhost:8000/admin/), and login with the super user credentials.  
Navigate to the [chunk servers page](http://localhost:8000/admin/master/chunkserver/), and add as many chunk servers as desired.


### 4. Start the chunk servers
```
cd apps/chunk-server
go run . --id {id} --port {port}
```
where the `id` and `port` must match the values from **step 3**.  
Repear this step for all servers marked as `is_active` in the master in **step 3**.

### 5. Initialize the web client
Create the `.env` file:
```
cd apps/web-client
echo "PUBLIC_SERVER_URL=http://localhost:8000" > .env
```
Run the web app:

```
pnpm install
pnpm dev
```

### 5. Upload a video
Use the `cli-client` to upload a `.mp4` file:
```
cd apps/cli-client
go run . --master-url http://localhost:8000 --file /path/to/video.mp4 --chunk-duration {chunk_duration}  --replicas {replication_factor} --title {title} --description {description}
```
where:
 - `chunk_duration` is the length of each chunk split in seconds except for possibly the last one.
 - `replication_factor` is the number of replicas for each chunk
 - `title` is a user-provided arbitrarily long string
 - `description` is a user-provided arbitrarily long string

### 6. Play a video
Open the [web app](http://localhost:5173/) and click any of the uploaded videos.
