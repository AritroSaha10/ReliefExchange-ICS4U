<div align="center">
<img width="300" src="https://github.com/AritroSaha10/ReliefExchange-ICS4U/assets/76495965/ed1cc35c-ec39-4066-b8a9-e0a5d53821d2" />
</div>

# ReliefExchange: ICS4U Final Project
ReliefExchange is a platform where **Generosity meets Community**. This website built by students at John Fraser Secondary School is meant for newcomers, immigrants, or any hard working individuals to find the essentials they need. In a society where individuals are struggling to be able to simply get the everyday items they need, and where waste is abundant, we built a platform where users can donate and trade for absolutely free! On the Relief Exchange platform, users sign up to create an account that allows them to either post an item they don’t need anymore, or browse through almost any category to find something for them. With how easy it is to simply throw something away, traditional donation centres are receving less and less items each year. By having an easy to use, online platform, we hope promote accessibility so users are more motivated to donate their items, instead of aiding to our ecological impact. To sign up now, visit https://reliefexchange.aritrosaha.ca/ 

# Features 
* User friendly interface that does not require the user to be well experienced with apps & Websites
* Ability to translate the website in over 6 different languages for immigrants and non-english speaking individuals
* Search to find Donation offers using filters including tags/categories, date uploaded, and more
* Upload a donation item with descriptions, images, location, and tags 
* Profile that tracks your donations, and users registration and sign-in date 

# Video Demonstration 
Our video demonstration can be seen at [this YouTube link](https://youtu.be/vHNYoDzi6TI).

# User Manual
Our user manual / quickstart guide is available in [this document](https://docs.google.com/document/d/1SwvbGomqzTCoZS3yOiGqkT-3oLW5rO9EociF8mV30MM/edit).

# Navigating Github
Our Github instructions can be seen in [this document](NAVIGATING_GITHUB.md).

# Built With 
- TypeScript
    - React
    - Next.js
- Golang
    - Gin
- Firebase
    - Firestore
    - Authentication
    - Storage
- Vercel (Frontend)
- Oracle Cloud Interface (Backend)
    - 1x VM w/ 4 Ampere ARM cores
    - Running as Docker container
    - NGINX (reverse proxy for container)
- GitHub Actions (CI/CD)
    - Both using GitHub-hosted servers and self-runners on OCI
- GitHub Container Registry (CI/CD)
- Testing Frameworks
    - React Testing Library
    - Jest
    - Cypress
    - `go-test`
    - `mockfs`
- Sentry (error logging for full stack)

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
 A Dockerfile and example Docker Compose file is available in the root directory of the backend folder. You can deploy this on any AMD64 or ARM64 machine. The built Docker image is available at [ghcr.io/aritrosaha10/relief-exchange-backend](https://github.com/aritrosaha10/ReliefExchange-ICS4U/pkgs/container/relief-exchange-backend). 
 
Our deployment runs this image as a singular Docker container on Oracle Cloud Infrastruction (OCI) behind a NGINX reverse proxy. It is highly recommended to use a reverse proxy such as Traefik or NGINX.

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

# User Testing Results
Our user testing results can be seen in [this document](USER_TESTING.md).

# Directory Outline
The description of each directory in our project can be seen in [this document](DIR_OUTLINE.md)

# Authors
This project was worked on by Aritro Saha, Joshua Chou, and Taha Khan as a part of our ICS4U0 (Computer Science) final project. 
