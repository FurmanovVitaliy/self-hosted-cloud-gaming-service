

export const event = (()=>{
const eventsStorage ={};

let _index = 0;
return {
    execute: (event, functionExecutor, order = undefined)=>{
        if(!eventsStorage[event]) eventsStorage[event] = {};
        let i = (order !== undefined ? order *1000000 : 0) + _index++;
        eventsStorage[event][i] = functionExecutor;
        return Object.freeze({
            terminate: ()=>{
             delete eventsStorage[event][i];
         }
        });
        
    },
    broadcast: (event,data)=>{
        if(!eventsStorage[event]) return;
        
        Object.keys(eventsStorage[event]).forEach((fnExe)=>{
            if (typeof eventsStorage[event][fnExe] === 'function') {
                eventsStorage[event][fnExe](data !== undefined ? data : {});
            }else{
               console.log("Not a function",eventsStorage[event][fnExe]);
            }

        });
    }
};
})();

export default event;