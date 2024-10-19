type FunctionExecutor = (data?: any) => void;
/**
 * A simple event bus implementation that allows for subscribing to and emitting events.
 */ 

class EventBus {
  storage: Record<string, Record<number, FunctionExecutor>> = {};
  index: number = 0;
  constructor() {}

 /**
   * Executes a function associated with an event, with optional execution once.
   * @param event The event name to listen for.
   * @param functionExecutor The function to execute when the event is triggered.
   * @param once Whether the function should only execute once.
   * @returns An object with a terminate method to stop listening for the event.
   */

  execute(event: string,functionExecutor: FunctionExecutor,once: boolean = false ): { terminate: () => void } {
    let idx = ++this.index; // Increment the global index for unique identifiers.
    if (!this.storage.hasOwnProperty(event)) {
      this.storage[event] = { [idx]: functionExecutor };
    } else {
      let idx = Object.values(this.storage[event]).length + 1;
      this.storage[event][idx] = functionExecutor;
    }
    const terminate = () => {
      if (!this.storage[event]) {
        return;
      }
      if (Object.keys(this.storage[event]).length === 0) {
        delete this.storage[event]; // Clean up the event slot if empty.
        return;
      }
      delete this.storage[event][idx];
    };

    if (once) {
      this.storage[event][idx] = (data: any) => {
        functionExecutor(data);
        terminate(); // Automatically terminate after executing once.
      };
    }
    return Object.freeze({ terminate });
  }

/**
   * Triggers all functions associated with the specified event.
   * @param event The event name to trigger.
   * @param data Optional data to pass to the event handlers.
   */

  notify(event: string, data?: any) {
    if (!this.storage[event]) {
      console.log(`No function waiting for event "${event}"`);
      return;
    }

    Object.values(this.storage[event]).forEach((func) => func(data));
     // Check for empty slots and clean up after notifications.
    if (Object.keys(this.storage[event]).length === 0) {
      delete this.storage[event];
    }
  }
}

const eventBus = new EventBus();
export default eventBus;
