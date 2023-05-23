# ReliefExchange (ICS4U Final Project)
ReliefExchange is a full-stack web platform where **generosity meets community**. We do this by connecting those in need with those who can give. This website, built by students at John Fraser Secondary School, is meant for newcomers, immigrants, or any hard working individuals to find the essentials they need. In a society where individuals are struggling to be able to simply get the everyday items they need, and where waste is abundant, we built a platform where users can donate and trade for absolutely free! On the Relief Exchange platform, users sign up to create an account that allows them to either post an item they don’t need anymore, or browse through almost any category to find something for them. With how easy it is to simply throw something away, traditional donation centres are receving less and less items each year. By having an easy to use, online platform, we hope promote accessinlity so users are more motivated to donate their items, instead of aiding to our ecological impact. To signup now visit https://reliefexchange-ics4u.vercel.app/

# Features 
* User friendly interface that does not require the user to be well experienced with apps & Websites
* Ability to translate the website in over 6 different languages for immigrants and non-english speaking individuals
* Search to find Donation offers using filters including tags/categories, date uploaded, and more
* Upload a donation item with descriptions, images, location, and tags 
* Profile that tracks your donations, and users registration and sign-in date 

# Built With 
* TypeScript
* Golang
* JavaScript

# Development Instructions
Make sure to change your working directory to the appropriate folder before running the commands.

## Frontend
- Install dependencies using `npm install`
- Copy `.env.template` into `.env.local`, and populate the environment variables with your values
- Run the development server using `npm run dev`
- Build the production code using `npm run build`
- Run the production server locally using `npm run start`
- Run tests using `npm run test` and `npm run cypress`

## Backend
- Install dependencies using `go mod download && go mod verify`
- Copy `.env.template` into `.env`, and populate the environment variables with your values
- Run the development server using `go run .`
- Build the production binary using `go build .`
- Run the production binary as you would any other executable binary file
- Run tests using `go test`

# Deployment Instructions

## Frontend
You have a few choices to deploy the frontend code.
- Deploy to Vercel (what we do) or Netlify (not tested)
- Deploy to another server (learn more about this [here](https://nextjs.org/docs/pages/building-your-application/deploying))
    - Please note that Static HTML Export is not possible for this project as we use Incremental Static Regeneration (ISR).

## Backend
A Dockerfile and example Docker Compose file is available in the root directory of the backend folder. You can deploy this on any AMD64 or ARM64 machine. The built Docker image is available at ghcr.io/aritrosaha10/reliefexchange-backend. Our deployment runs this image as a singular Docker container on Oracle Cloud Infrastruction (OCI), using the Ampere ARM cores. It is highly recommended to use a reverse proxy 
such as Traefik or NGINX.

## Firebase
- Initialize a new Firebase project
- Add your own configuration data to the appropiate frontend and backend
- Set up:
    - Firestore
        - Copy-paste the security rules found [here](firebase/rules/firestore.cel)
        - Create three collections:
            - `config`
                - This should have one document called `bans`
                    - This should have an array with the key of `users`
            - `donations`
            - `users`
    - Authentication
        - Add Google as a sign-in provider
    - Storage
        - Copy-paste the rules found [here](firebase/rules/storage.cel)
        - Create a folder called `donations`
