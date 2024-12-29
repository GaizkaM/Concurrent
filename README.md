# Concurrent Programming Projects

This repository contains three different projects developed for the **Concurrent Programming** course. Each project explores various concurrency problems using different synchronization mechanisms and programming languages.

---

## 📂 Projects Overview

### 1. **Santa Claus and the Elves Problem** 🧑‍🎄🎁
- **Description**: Simulates a concurrent system where Santa Claus is assisted by elves and reindeer. Each elf builds toys and periodically needs Santa's help. Once all toys are completed, Santa prepares the sleigh with the reindeer and finishes his task.
- **Key Features**:
  - Elves work in groups of three to request Santa's assistance.
  - Santa rests when not needed and wakes up when summoned.
  - Reindeer arrive asynchronously and wait to be hitched to the sleigh.
  - The simulation concludes when all toys are made, and the reindeer are ready to depart.
- **Technologies**: Implemented in **Java** or **Python** using semaphores for synchronization.
- **Highlights**:
  - Processes: `Santa`, `Elf`, `Reindeer`.
  - Key synchronization via semaphores to handle waiting rooms and task queues.

---

### 2. **The Soda Machine Problem** 🥤
- **Description**: A producer-consumer problem where a soda machine is shared between random customers and suppliers. Customers consume sodas while suppliers refill the machine, all managed via concurrent processes.
- **Key Features**:
  - Customers choose a random number of sodas to consume.
  - Suppliers refill the machine until all customers finish.
  - The machine acts as a monitor ensuring proper synchronization.
  - Random delays simulate real-time interactions.
- **Technologies**: Implemented in **Ada** using protected objects for synchronization.
- **Highlights**:
  - Processes: `Customer`, `Supplier`.
  - The soda machine is a shared resource managed as a monitor.
  - Messages for synchronization include state transitions like "machine is empty" or "machine refilled".

---

### 3. **The Illegal Tobacco Shop Problem** 🚬
- **Description**: Based on the smokers problem, this simulation features a tobacco shop, smokers with specific needs (matches or tobacco), and a whistleblower. The shopkeeper provides resources to smokers until interrupted by the police.
- **Key Features**:
  - Smokers with matches wait for tobacco, and smokers with tobacco wait for matches.
  - The shopkeeper serves requests in a round-robin fashion.
  - A whistleblower notifies all processes to stop when the simulation ends.
- **Technologies**: Implemented in **Go** using **RabbitMQ** for message passing.
- **Highlights**:
  - Processes: `Shopkeeper`, `SmokerWithMatches`, `SmokerWithTobacco`, `Whistleblower`.
  - RabbitMQ’s `Fanout Exchange` ensures all processes receive termination signals.

---

## 📝 Author
- **Gaizka Medina Gordo**
