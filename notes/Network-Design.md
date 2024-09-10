# Poker Game api

## Js types

```JS
/** Poker Cards {A23456789XJQK}{♣♦♥♠}
 * @typedef {'♣'|'♦'|'♥'|'♠'} suit
 * @typedef {'A'|'2'|'3'|'4'|'5'|'6'|'7'|'8'|'9'|'X'|'J'|'Q'|'K'} rank
 * @typedef {`  `|`${rank}${suit}`} card
 */
```

```JS
/** Poker Game Object
 *
 * @typedef {'BuyIn'|'Check'|'Bet'|'Winner'} phase
 * @typedef {'Check'|'Fold'|'Call'|'Raise'} action
 * @typedef {number} cents
 *
 * @typedef {{ // Board
 *  pot    : cents,
 *  minBet : cents,
 *  round  : number,
 *  cards  : [card, card, card, card, card],
 *
 *  deck   : undefined,
 * }} board
 *
 * @typedef {{ // Player
 *  id        : string,
 *  name      : string,
 *
 *  cards     : [card, card],
 *  score     : number,
 *  handTypes : string[],
 *
 *  boughtIn  : boolean,
 *  folded    : boolean,
 *  action    : action,
 *
 *  bet       : cents,
 *  bank      : cents,
 * }} player
 *
 * @typedef {{ // Game
 *  id          : string,
 *  board       : board,
 *  players     : player[], // [1..8]
 *  playerTurnI : number,
 *  phase       : phase,
 * }} game
 *
 */
```

## Routes

| Method | Route               | Description                       | Response                                    |
| ------ | ------------------- | --------------------------------- | ------------------------------------------- |
| GET    | /poker              | Game Webpage                      | `html`                                      |
| GET    | /poker/{gameId}     | Game Webpage after joining a game | `html`                                      |
| POST   | /api/poker/create   | Create a new game                 | `=>{ game:game } `<br>`\| { error:string }` |
| WS     | /api/poker/{gameId} | Connect to a game                 | `tcp`                                       |

## Socket Messages

### Client To Server

- ```JS
  // Client
  async ( { type:'JoinAsSpectator' }
  // Server
    ) =>  { type:'JoinAsSpectator' }
        | { error:string }
  ```

- ```JS
  // Client
  async ( { type:'NewPlayer', name:string }
  // Server
    ) =>  { type:'NewPlayer', you:player, token:string }
        | { error:string }
  ```

- ```JS
  // Client
  async ( { type:'JoinAsPlayer', playerId:string, token:string }
  // Server
    ) =>  { type:'JoinAsPlayer' }
        | { error:'NotFound'|'Auth'|'Full'|'Present' }
  ```

- ```JS
  // Client
  async ( { type:'LeaveAsPlayer', playerId:string }
  // Server
    ) =>  { type:'LeaveAsPlayer' }
        | { error:'Join' }
  ```

- ```JS
  // Client
  async ( { type:'GameState' }
  // Server
    ) =>  { type:'GameState', game:game }
        | { error:'Join' }
  ```

- ```JS
  // Client
  async ( { type:'BoardState' }
  // Server
    ) =>  { type:'BoardState', board:board }
        | { error:'Join' }
  ```

- ```JS
  // Client
  async ( { type:'PlayerState', playerId:string }
  // Server
    ) =>  { type:'PlayerState', player:player }
        | { error:'NotFound' }
  ```

- ```JS
  //Client
  async ( { type:'BuyIn' }
  // Server
    ) =>  { type:'BuyIn' }
        | { error:'Join'|'Funds' }
  ```

- ```JS
  //Client
  async ( { type:'Action', action:'Check'|'Fold'|'Call' }
        | { type:'Action', action:'Raise', raiseTo:number }
  // Server
    ) =>  { type:'Action' }
        | { error:'Join'|'Phase'|'BuyIn'|'Folded' }
  ```

### Server to Client

- ```JS
  // Server
  async ( { type:'PlayerJoin', player:player, seatedAfterPlayerId:string }
  // Client
    ) =>  { type:'PlayerJoin' }
  ```

- ```JS
  // Server
  async ( { type:'PlayerLeave', player:player }
  // Client
    ) =>  { type:'PlayerLeave' }
  ```

- ```JS
  // Server
  async ( { type:'PlayerBuyIn', player:player }
  // Client
    ) =>  { type:'PlayerBuyIn' }
  ```

- ```JS
  // Server
  async ( { type:'PlayerAction', player:player }
  // Client
    ) =>  { type:'PlayerAction' }
  ```

- ```JS
  // Server
  async ( { type:'RoundStart', you:player }
  // Client
    ) =>  { type:'RoundStart' }
  ```

- ```JS
  // Server
  async ( { type:'RoundEnd', winner:player }
  // Client
    ) =>  { type:'RoundEnd' }
  ```

- ```JS
  // Server
  async ( { type:'PhaseStart', phase:phase }
  // Client
    ) =>  { type:'PhaseStart' }
  ```

- ```JS
  // Server
  async ( { type:'BoardCheck', board:board }
  // Client
    ) =>  { type:'BoardCheck' }
  ```
