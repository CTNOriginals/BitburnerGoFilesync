export function main(ns: NS) {
	var x = 0;
	ns.tprint("hello" + x)

	var host = (ns.args.length > 0) ? ns.args[0] as string : "home"

	var secLevel = ns.getServerSecurityLevel(host)
	var maxSec = ns.getServerMinSecurityLevel(host)

	var money = ns.getServerMoneyAvailable(host)

	var growth = ns.getServerGrowth(host)
	var maxMoney = ns.getServerMaxMoney(host)

	ns.tprintf("money: %f\nsec: %f.2/%f.2\ngrowth: %f.2/%f.2", money, secLevel, maxSec, growth, maxMoney)
}
