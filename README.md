# Image Labeling System Based On Microservices Architecture

## Overview

The main challenge in training supervised learning models is having a sufficient amount of accurately labeled data. Acquiring a large volume of correctly labeled images for model training is far from trivial, as it requires significant time and effort, especially with limited personnel. To address this issue, we designed a system that allows clients to upload a set of images for classification. These images are first uploaded to the system and then distributed to various classifiers, who are tasked with assigning labels to each image. Once the images are fully classified, they are returned to the client with their associated labels.

### Scope

The scope of this project includes:
- Investigating microservices architecture from theoretical and practical perspectives.
- Developing a system capable of efficient orchestration and communication among microservices.
- Ensuring the system's scalability to handle dynamic workloads.

More details on the [Techincal Report](technical_report_rcd.pdf)

# System Overview

The Image Labeling System is designed around two core components: the Task Manager and the Subtask Manager, each playing a pivotal role in managing the workflow of tagging images with labels for training machine learning models.

![system](https://github.com/alecava41/image-labeling-system/assets/76614857/3e21e758-4989-4157-afe5-754f177a7e5a)


### Task Manager

The task manager is responsible for overseeing the submission of a group of images to be tagged. Additionally, it is tasked with assigning a label to a specific subtask upon its completion.

From a functional point of view, here is an overview of how the service behaves. Clients send a potentially huge block of images. Each block of images to be tagged corresponds to a task to be performed by the system. The block of images must be sent along with a set of labels: this one will be used by operators to tag individual images.

From the point of view of the implementation, upon reception of a request from a client, the service must process all the images following these steps:
1. Initially, it will create a task, intended as the task of tagging those images with the provided labels;
2. It must validate that a particular uploaded file is an image;
3. It must upload the images into the storage (MinIO bucket). 
4. On upload completion, it must populate the previously-created task with a new subtask, representing the duty of tagging the just uploaded image;
5. Lastly, once the creation of the subtask (on the database) has been successful, the service must publish a new event into the message broker, stating that a new image is available and ready to be tagged by operators.

Here is a graphical overview about the flow of such a service.

![task-manager](https://github.com/alecava41/rcd-project/assets/76614857/c9bf01bd-5667-4cec-bc08-3337ab3c7f25)


### Subtask Manager

The subtask manager is responsible for providing the images (subtasks) requested by the operators. Once they have classified them, it has the task of storing, for each subtask, the various labels received. As soon as the necessary number of labels is reached, the subtask will be considered completed, therefore the task manager will be notified, along with the information related to the various labels assigned to that specific image.

An operator initiates the process by requesting a group of subtasks, each representing an image to be classified. These subtasks are dispatched to the requesting operator and include both the ID of the image and a selection of possible labels that can be assigned to it. The operator then assigns an appropriate label to each image, and these classification outcomes are communicated back to the subtask manager.

Here's a more detailed description of how the workflow is implemented:
1. The subtask-manager creates subtasks in response to notifications from the task manager triggered by the upload of individual images. A subtask is defined by the image's ID, a set of assignable labels for that image, and a list of labels that will be assigned by operators.
2. When an operator requests a set of images, the service provides a set of images (subtasks) for completion.
3. After the operator submits the classified images, the service updates each subtask by adding the new label provided by the operator.
4. Once a subtask accumulates a sufficient number of labels, signaling that the classification process is complete, the service informs the task manager of this achievement, indicating that the subtask has met its required label count.

Below is presented a graphical overview of flow when an image has been classified by an operator and the result is sent back to the subtask-manager service.

![subtask-manager](https://github.com/alecava41/rcd-project/assets/76614857/25b7616b-1f2e-4164-bdab-cd28f83b81c9)

## Experimental Implementation

### Main Technologies

- **Programming Language:** Go, for its concurrency support and efficient resource management.
- **Web Framework:** Gin, for creating RESTful APIs.
- **Container Orchestration:** Kubernetes, for automating deployment, scaling, and management.
- **Message Broker:** RabbitMQ, for asynchronous communication between microservices.
- **Database:** MongoDB, for storing tasks and subtasks.
- **Image Storage:** MinIO, for storing images to be tagged.

### Design of Evaluation Experiments

Our evaluation focuses on the system's ability to scale in response to varying workloads. We simulate client and operator behavior using custom Go programs, assessing the system's performance and scalability through a series of load tests.

## Results and Discussion

The experimental results demonstrate the system's ability to scale effectively in response to workload changes. The use of Kubernetes and load balancers allows the system to handle increased workloads by deploying additional replicas and to conserve resources by scaling down when demand decreases.

## Limitations

- Lack of authentication and request constraints.
- External services (e.g., MongoDB, RabbitMQ, MinIO) were not configured to autoscale, which could be a limitation in real-world scenarios.

## Conclusion

The project has provided valuable insights into the design, implementation, and scaling of microservices-based systems. Through practical experimentation, we have validated the benefits of microservices architecture and identified areas for future improvement.
