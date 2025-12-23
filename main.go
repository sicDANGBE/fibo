package main

import (
	"fmt"
	"math/big"
)

// FibGenerator produit la suite de Fibonacci via un channel.
// S'arrête si la RAM estimée des deux derniers nombres dépasse limitBytes.
func FibGenerator(limitBytes uint64) <-chan *big.Int {
	ch := make(chan *big.Int)

	go func() {
		defer close(ch)
		// Initialisation : F(0)=0, F(1)=1
		a := big.NewInt(0)
		b := big.NewInt(1)

		for {
			// Calcul de la taille approximative en RAM des deux termes
			// big.Int stocke les données dans un slice de 'Word' (uint sur 64 bits)
			// On compte environ 8 octets par Word + le overhead de la structure.
			sizeA := uint64(len(a.Bits())) * 8
			sizeB := uint64(len(b.Bits())) * 8

			if sizeA+sizeB > limitBytes {
				fmt.Printf("\n[Limite de %d Go atteinte]\n", limitBytes/1024/1024/1024)
				return
			}

			// On envoie une copie pour éviter les effets de bord si l'appelant modifie la valeur
			val := new(big.Int).Set(a)
			ch <- val

			// Fibonacci : a, b = b, a+b
			// On utilise Add pour additionner b à a, puis on swap.
			a.Add(a, b)
			a, b = b, a
		}
	}()

	return ch
}

func main() {
	const maxRAM = 5 * 1024 * 1024 * 1024 // 5 Go
	gen := FibGenerator(maxRAM)

	count := 0
	for f := range gen {
		count++
		// Pour l'exemple, on affiche tous les 100 000 termes
		// car l'affichage console est très lent pour de gros nombres.
		if count%100000 == 0 {
			fmt.Printf("Terme n°%d calculé (Taille actuelle : ~%d MB)\n", count, len(f.Bits())*8/1024/1024)
		}
	}
}
