const message ={};

export const  msg = (
  contentType,
  content,
  UUID = undefined,
  username = undefined
) => {
  message.content_type = contentType;
  message.content = content;
  if (typeof content === "object") {
    message.content = btoa(JSON.stringify(message.content));
  }
  message.UUID = UUID;
  message.username = username;
  return JSON.stringify(message);
};
export default msg;