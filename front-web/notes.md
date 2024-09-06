# Poker Game api

## Js types

```JS
/** Poker Cards {A23456789XJQK}{♣♦♥♠}
 * @typedef {'♣'|'♦'|'♥'|'♠'} suit
 * @typedef {'A'|'2'|'3'|'4'|'5'|'6'|'7'|'8'|'9'|'X'|'J'|'Q'|'K'} rank
 * @typedef {`  `|`${rank}${suit}`} card
 */

/** Poker Game Object
 *
 * @typedef {{
 *  pot   : number,
 *  round : number,
 *  cards : [card, card, card, card, card],
 * }} board
 *
 * @typedef {{
 *  id       : string,
 *  name     : string,
 * 
 *  cards    : [card, card],
 *  score    : number,
 *  handType : string[],
 * 
 *  bet      : number,
 *  bank     : number,
 * }} player
 *
 * @typedef {{
 *  id          : string,
 *  board       : board,
 *  players     : player[], // [1..8]
 *  playerTurnI : number,
 * }} game
 *
 */
```

## Routes

| Method | Route                   | Description                           | Response                                 |
| ------ | ----------------------- | ------------------------------------- | ---------------------------------------- |
| GET    | poker                   | Game Webpage                          | `html`                                   |
| GET    | poker/{gameId}          | Game Webpage after joining a game     | `html`                                   |
| POST   | poker/api/game/create   | Create a new game                     | `=>{ game:game   } `<br>`\| { error:string }` |
| POST   | poker/api/player/create | Create a new player (with auth token) | `=>{ player:player } `<br>`\| { error:string }`                                         |
| WS     | poker/api/{gameId}      | Connect to a game                     | `tcp`                                        |

## Ws messages

### Client To Server

- ```JS
  // Client
  async ( { type:'Connect', as:`${playerId}`, auth:string }
        | { type:'Connect', as:'Spectator'  }
  // Server
    ) =>  { type:'Connect', accept:true }
        | { error:'Present'|'Auth' }
  ```

- ```JS
  // Client
  async ( { type:'GetState' }
  // Server
    ) =>  { type:'GameState', gameState:game }
        | { error:'Join' }
  ```

- ```JS
  //Client
  async ( { type:'Move', action:'Check'|'Fold'|'Call' }
        | { type:'Move', action:'Raise', raiseTo:number }
  // Server
    ) =>  { type:'Move', gameState:{...} }
        | { error:'Join'|'Turn'|'NoMoreCheck'|'Funds' }
  ```

### Server to Client

- ```JS
  // Server
  async ( { type:'RoundStart', you:player }
  // Client
    ) =>  { type:'RoundStart' }
  ```

- ```JS
  // Server
  async ( { type:'Move', player:player, nextPlayer:player, board?:board }
  // Client
    ) =>  { type:'Move' }
  ```

- ```JS
  // Server
  async ( { type:'RoundEnd', game:game, winner:player }
  // Client
    ) =>  { type:'RoundEnd' }
  ```
