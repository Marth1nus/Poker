/** Poker Cards {A23456789XJQK}{♣♦♥♠}
 * @typedef {'♣'|'♦'|'♥'|'♠'} suit
 * @typedef {'A'|'2'|'3'|'4'|'5'|'6'|'7'|'8'|'9'|'X'|'J'|'Q'|'K'} rank
 * @typedef {`${rank}${suit}`} card
 */
/** Poker Game Object
 * @typedef {{pot:number, cards:card[]}} board
 * @typedef {{id:string, hand:card[], name:string, bet:number, bank:number}} player
 */

document.addEventListener("DOMContentLoaded", onDOMContentLoaded, { once: true })

/** @param {Event} event  */
function onDOMContentLoaded(event) {
  /** @type {HTMLFormElement} */
  const form = document.querySelector("#poker-game .player form")
  form.addEventListener("submit", onPlayerFormSubmit)
}

/**
 * @param {SubmitEvent} event
 * @this {HTMLFormElement}
 */
async function onPlayerFormSubmit(event) {
  event.preventDefault()

  /** @type {{action:string, raiseTo:number}} */
  const formData = Object.fromEntries(new FormData(this).entries()),
    { action, raiseTo } = formData
  console.info("Player Input", formData)
  if (!action) {
    console.error("No Action")
    return
  }

  const response = await fetch(`/api/{gameId}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ action, raiseTo: action === "Raise" ? raiseTo : undefined }),
  })

  if (!response.ok) {
    console.error("Response Fail", response)
    return
  }

  /** @type {{prevPlayer:player, currPlayer:player, currBoard:board}} */
  const { prevPlayer, currPlayer, currBoard } = await response.json()

  // Update player seats
  for (const player of [prevPlayer, currPlayer]) {
    /** @type {HTMLElement} */
    const hand = document.querySelector(`#poker-game .seat[data-playerId="${player.id}"] .hand`)
    hand.classList.toggle("playing", player === currPlayer)

    /** @type {HTMLElement[]} */
    const [card0, card1, name, betTitle, bet, bankTitle, bank] = hand.children
    if (player.hand?.length === 2)
      [
        [card0.dataset.rank, card0.dataset.suit], //
        [card1.dataset.rank, card1.dataset.suit],
      ] = player.hand
    name.innerText = player.name
    bet.innerText = player.bet
    bank.innerText = player.bank
  }
}
