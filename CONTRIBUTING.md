# Contribution Guide
Thanks for your great contribution!!

## Principles
- Let's be creative and collaborative👶
- Let's welcome people from different backgrounds👶
- This project follows [Go Community Code of Conduct](https://go.dev/conduct)👶

## Issues
- When you report a new issue, please clearly mention the following three points at least 🎉
  - what happened/is happening
  - why that is a problem
  - solution hypothesis

## Pull Request
- If possible, please report an issue at first😉
- **But the project always welcome your collaborative pull request!**
- **Please fork the develop branch to create PR🎉**

## Before Pushing...
### Architecture
- sqloth is on DDD architecture.
  - The main logic is in ```domain```, especially, in ```model``` and ```dservice```(=domain service).
  - The ```driver``` layer processes (given) external info. The interface is in ```domain/driver``` and the implementation is in ```driver```.
    - Currently, the layer has only one module, ```file_driver```. This is because any inputs other than a schema file are NOT supposed.
  - Please combine them in the ```usecase``` layer to implement features.

### Unit Test
- Please write unit tests using go test
- If you define new interfaces, please generate mock by editing ```make gen-mock```commands and running it.
- At last, please check if your change pass all the tests by running the below locally.
```bash
make test
```
