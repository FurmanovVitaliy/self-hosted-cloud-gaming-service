import constants from "./constants";


type LogLevel = 'info' | 'warn' | 'error' | 'debug';
class Logger {
  private static enabled = constants.log.enabled;

  private static getColor(level: LogLevel): string {
    if (level === 'info') return '\x1b[94m';
    if (level === 'warn') return '\x1b[93m';
    if (level === 'error') return '\x1b[91m';
    if (level === 'debug') return '\x1b[90m';
    return '\x1b[0m';  // Reset
  }

  private static getTimestamp(): string {
      return new Date().toLocaleTimeString();
  }

  private static getStackTrace(): string {
      const err = new Error();
      const stack = err.stack || '';
      const stackLines = stack.split('\n');
      // Индекс 4, потому что первые три строки стека - это текущий вызов getStackTrace, метод логгирования и конструктор Error
      const callerLine = stackLines[4] || '';
      const match = callerLine.match(/at (.+):(\d+):(\d+)/);
      if (match) {
          return `${match[1]}:${match[2]}`;
      }
      return 'unknown location';
  }
  private static getLogLevel(leves: LogLevel) : number {
    switch (leves) {
      case 'debug':
        return 4;
      case 'info':
        return 3;
      case 'warn':
        return 2;
      case 'error':
        return 1;
      default:
        return 0;
    }
  }

  private static log(level: LogLevel, message: any ): void {
      if (!Logger.enabled) {
          return;
      }
      if (Logger.getLogLevel(level) > constants.log.level) {
        return;
      }
      const color = Logger.getColor(level);
      const time = Logger.getTimestamp();
      const location = Logger.getStackTrace();
      console.log(color)
     if (typeof message === 'object') {
      console.log(`${color}${time} [${level.toUpperCase()}] ${JSON.stringify(message)} (at ${(location)})\x1b[0m`);
      return;
    }
      console.log(`${color}${time} [${level.toUpperCase()}] ${message} (at ${(location)})\x1b[0m`);
  }

  static debug(message: any): void {
    Logger.log('debug', message);
  }
  
  static info(message: any): void {
      Logger.log('info', message);
    
  }

  static warn(message: any): void {
      Logger.log('warn', message);
  }

  static error(message: any): void {
      Logger.log('error', message);
  }
}
const log = Logger;
export default log;