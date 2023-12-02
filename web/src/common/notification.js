export const notification = (()=>{
    const storage = {}
    let index = 0;
    return {
        execute: (event, functionExecutor, isOnce = false)=>{ 
            let once = isOnce ? 0 : 1;
        index++
           
           if (!storage.hasOwnProperty(event)) {
               storage[event] = { [(index*once)]: functionExecutor };
            } else {
                storage[event][index*once] = functionExecutor;
            }
            
            function terminate(){delete storage[event][index]}
            return Object.freeze({terminate})
        },
        
        broadcast: (event,data)=>{
            if(!storage[event]) return console.log(`No function witch waiting for event "${event},or the order is incorrect"`);
            Object.keys(storage[event]).forEach((functionExecutor)=>{
                
                if (typeof storage[event][functionExecutor] !== 'function') return console.log(`Function "${functionExecutor}" is not a function.\n Corect function exemple: ' console.log ' `);
                storage[event][functionExecutor](data !== undefined ? data : {})
             }
            )
        }
    }
}
)()

export default notification;