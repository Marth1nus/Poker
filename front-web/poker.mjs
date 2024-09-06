/** Poker Cards {A23456789XJQK}{♣♦♥♠}
 * @typedef {'♣'|'♦'|'♥'|'♠'} suit
 * @typedef {'A'|'2'|'3'|'4'|'5'|'6'|'7'|'8'|'9'|'X'|'J'|'Q'|'K'} rank
 * @typedef {`  `|`${rank}${suit}`} card
 */

/** Poker Game Object
 *
 * @typedef {'Fold'|'Call'|'Raise'} action
 *
 * @typedef {{
 *  id    : string,
 *  cards : card[],
 *  name  : string,
 *  bet   : number,
 *  bank  : number,
 * }} player
 *
 * @typedef {{
 *  id          : string,
 *  board       : { pot:number, cards:card[] },
 *  players     : player[],
 *  playerTurnI : number
 * }} game
 *
 */

/**
 * 
 */

const pokerGame = {
  get gameDiv() {
    return window.location.pathname.match(/game\/(.+?)\/?/g)?.at(1)
  },
  get playerId() {
    return
  }
}

/** Create a new Game
 * POST `api/game`
 * @returns {Promise<{game:game}|{error:string}>}
 */
async function api_POST_game() {
  const url = `api/game`
  const response = await fetch(url, { method: 'POST' })
  /** @type {{game:game, error?:string}} */
  const { game, error = response.statusText } = await response.json()
  switch (response.status) {
    case 201 /* Created */:
      return { game }
    default:
      return { error }
  }
}

/** Get Up to data game state
 *
 * GET `api/game/${gameId}&modifiedOnly=${!!modifiedOnly}`
 * @param {string} gameId
 * @returns {Promise<
 *  {modifiedOnly:true, modified:true, game:game} |
 *  {modifiedOnly:true, modified:false} |
 *  {modifiedOnly:false, modified:boolean, game:game} |
 *  {error:string}>}
 */
async function api_GET_game_gameId(gameId, modifiedOnly = true) {
  const url = `api/game/${gameId}&modifiedOnly=${!!modifiedOnly}`
  const response = await fetch(url, { method: 'GET' })
  /**
   * @type {{game:game, modified:boolean, error?:string}} */
  const { game, modified, error = response.statusText } = await response.json()
  switch (response.status) {
    case 304 /* Not Modified */:
    case 200:
      return modifiedOnly //
        ? modified
          ? { modifiedOnly, modified, game }
          : { modifiedOnly, modified }
        : { modifiedOnly, modified, game }
    case 404 /* Not Found */:
    default:
      return { error }
  }
}

/** Post player move
 *
 * POST `api/game/${gameId}`
 * @param {string} gameId
 * @param {{playerId:string, action:action, raiseTo?:number}} RequestBody
 * @return {Promise<{game:game}|{error:string}>}
 */
async function api_POST_game_gameId(gameId, { playerId, action, raiseTo }) {
  const url = `api/game/${gameId}`
  const response = await fetch(url, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ playerId, action, raiseTo }),
  })
  /** @type {{game:game, error?:string}} */
  const { game, error = response.statusText } = await response.json()
  switch (response.status) {
    case 202 /* Accepted */:
      return { game }
    case 417 /* Expectation Failed */:
    default:
      return { error }
  }
}

/** Apply a game object to the DOM
 * @param {Element} gameDiv
 * @param {game} game
 */
async function applyGameState(gameDiv, game) {
  if (!(gameDiv instanceof HTMLElement)) throw TypeError('gameDiv is not a html element')

  gameDiv.id = game.id
  gameDiv.dataset.playerId = game.players.at(0)?.id

  const pot = gameDiv.querySelector('.board p')
  if (pot) pot.textContent = (game.board.pot / 100).toFixed(2)

  gameDiv //
    .querySelectorAll('.board .card')
    .forEach((card, i) => applyCardState(card, game.board.cards[i]))

  gameDiv //
    .querySelectorAll('.seat')
    .forEach((seat, i) => {
      seat.classList.toggle('playing', i === game.playerTurnI)

      /** @type {player|undefined} */
      const player = game.players[i]

      if (player) seat.id = player.id
      else seat.toggleAttribute('id', false)

      seat //
        .querySelectorAll('.hand .card')
        .forEach((card, i) => applyCardState(card, player?.cards[i]))

      const name = seat.querySelector('.hand b')
      if (name) name.textContent = player.name

      const [bet, bank] = seat.querySelectorAll('.hand p')
      if (bet) bet.textContent = ((player?.bet || 0) / 100).toFixed(2)
      if (bank) bank.textContent = ((player?.bank || 0) / 100).toFixed(2)

      const form = seat.querySelector('form')
      if (form && player) {
        form.dataset.gameId = game.id
        form.dataset.playerId = player.id
      }
    })
}

/** Apply a card string to a div
 * @param {Element} div
 * @param {card} [card]
 * @throws {TypeError}
 */
function applyCardState(div, card) {
  if (!(div instanceof HTMLElement)) throw TypeError('Card can only be applied to div')
  ;[div.dataset.rank, div.dataset.suit] = card || '  '
}

document.addEventListener('DOMContentLoaded', onDOMContentLoaded, { once: true })

/** Event: Dom Content Loaded
 * @param {Event} event
 */
async function onDOMContentLoaded(event) {
  try {
    const gameDiv = document.querySelector('.poker-game:not([id])')
    if (!gameDiv) throw new Error('no game div')

    /** @type {HTMLFormElement|null} */
    const form = gameDiv.querySelector('.poker-game .player form')
    if (!form) throw new Error('no player form')
    form.addEventListener('submit', onPlayerFormSubmit)

    /**
     * @param {MouseEvent} event
     * @this {HTMLInputElement}
     */
    const setValue = function (event) {
      const inputButton = this
      const form = inputButton.parentElement
      const inputAction = form?.querySelector('input[name="action"]')
      if (!(inputAction instanceof HTMLInputElement)) throw new TypeError()
      inputAction.value = inputButton.value
    }
    for (const input of form.querySelectorAll('input[type="submit"]')) {
      if (!(input instanceof HTMLInputElement)) throw new TypeError()
      input.addEventListener('click', setValue)
    }

    const gameId = pokerGame.gameDiv
    if (!gameId) throw new Error('No gameId in url')

    const responseBody = await api_GET_game_gameId(gameId, true)
    if ('error' in responseBody) throw new Error(responseBody.error)
    else if (responseBody.modified) {
      await applyGameState(gameDiv, responseBody.game)
      console.info('Game loaded', responseBody.game)
    }

    const updateInterval = setInterval(pullUpdate, 2000)
    async function pullUpdate() {
      const gameId = pokerGame.gameDiv
      if (!gameId) throw new Error('No gameId in url')
      const responseBody = await api_GET_game_gameId(gameId, true)
      if ('error' in responseBody) throw new Error(responseBody.error)
      const gameDiv = document.getElementById(gameId)
      if (!gameDiv) throw new Error(`Game has no div: ${gameId}`)
      if (!responseBody.modified) return
      return applyGameState(gameDiv, responseBody.game)
      clearInterval(updateInterval)
    }
  } catch (err) {
    // Serverless behavior
    // alert('failed to load the game')
    const playingInterval = setInterval(() => {
      const [player, currPlaying, nextPlaying] = ['.seat.player', '.seat.playing', '.seat.playing+.seat'] //
        .map(document.querySelector, document)
      ;(currPlaying || player)?.classList.toggle('playing', false)
      ;(nextPlaying || player)?.classList.toggle('playing', true)
      if (!nextPlaying) clearInterval(playingInterval)
    }, 1000)
    throw err
  }
}

/** Event: Player Input Submit
 * @param {SubmitEvent} event
 * @this {HTMLFormElement}
 */
async function onPlayerFormSubmit(event) {
  event.preventDefault()

  const formData = Object.fromEntries(new FormData(this).entries())
  const { action, raiseTo } = formData
  console.info('Player Input', formData)

  const { gameId, playerId } = this.dataset
  if (!gameId) throw new Error(`Blank gameId ${gameId}`)
  if (!playerId) throw new Error(`Blank playerId ${playerId}`)

  const responseBody = await (async () => {
    switch (action) {
      case 'Fold':
      case 'Call':
        return api_POST_game_gameId(gameId, { playerId, action })
      case 'Raise':
        return api_POST_game_gameId(gameId, { playerId, action, raiseTo: Number(raiseTo) })
      default:
        throw Error(`Unknown action, ${action}`)
    }
  })()
  if ('error' in responseBody) throw new Error(responseBody.error)
  const gameDiv = document.querySelector('.poker-game')
  if (gameDiv) applyGameState(gameDiv, responseBody.game)
}
