/** {A23456789XJQK}{♣♦♥♠}
 * @typedef {'♣'|'♦'|'♥'|'♠'} suit
 * @typedef {'A'|'2'|'3'|'4'|'5'|'6'|'7'|'8'|'9'|'X'|'J'|'Q'|'K'} rank
 * @typedef {`${rank}${suit}`} card
 * @typedef {HTMLDivElement} card
 */

/** @type    {rank[]} */ export const ranks = Array.from("A23456789XJQK")
/** @type    {suit[]} */ export const suits = Array.from("♣♦♥♠")
/** @returns {card[]} */ export const allCards = () => ranks.flatMap(rank => suits.map(suit => `${rank}${suit}`))
/** @returns {card[]} */ export const newDeck = () => shuffle(allCards())



/** Shuffle items in place
 * @template T
 * @param {T[]} items
 * @returns {T[]} items
 */
export function shuffle(items) {
  for (let i = items.length - 1; 0 < i; i--) {
    const j = Math.floor(Math.random() * (i + 1))
    ;[items[i], items[j]] = [items[j], items[i]]
  }
  return items
}
