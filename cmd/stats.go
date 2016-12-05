package cmd

import (
	"fmt"
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
	"os"
)

var statsintreestr string
var statsoutfile string
var statintrees chan tree.Trees
var statsout *os.File

// statsCmd represents the stats command
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Displays different statistics about the tree",
	Long: `Displays different statistics about the tree

For exemple:
- Edge informations
- Node informations
- Tips informations

`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		var nbtrees int = 0
		statintrees = make(chan tree.Trees, 15)
		/* Read ref tree(s) */
		go func() {
			if nbtrees, err = utils.ReadCompTrees(statsintreestr, statintrees); err != nil {
				io.ExitWithMessage(err)
			}
		}()
		statsout = openWriteFile(statsoutfile)
	},
	Run: func(cmd *cobra.Command, args []string) {
		/* Dividing trees */
		statsout.WriteString("tree\tnodes\ttips\tedges\tmeanbrlen\tsumbrlen\tmeansupport\tmediansupport\trooted\n")
		for statsintree := range statintrees {
			statsintree.Tree.ComputeDepths()
			statsout.WriteString(fmt.Sprintf("%d", statsintree.Id))
			statsout.WriteString(fmt.Sprintf("\t%d", len(statsintree.Tree.Nodes())))
			statsout.WriteString(fmt.Sprintf("\t%d", len(statsintree.Tree.Tips())))
			statsout.WriteString(fmt.Sprintf("\t%d", len(statsintree.Tree.Edges())))
			statsout.WriteString(fmt.Sprintf("\t%.4f", statsintree.Tree.MeanBrLength()))
			statsout.WriteString(fmt.Sprintf("\t%.4f", statsintree.Tree.SumBranchLengths()))
			statsout.WriteString(fmt.Sprintf("\t%.4f", statsintree.Tree.MeanSupport()))
			statsout.WriteString(fmt.Sprintf("\t%.4f", statsintree.Tree.MedianSupport()))
			if statsintree.Tree.Rooted() {
				statsout.WriteString(fmt.Sprintf("\trooted\n"))
			} else {
				statsout.WriteString(fmt.Sprintf("\tunrooted\n"))
			}
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		statsout.Close()
	},
}

func init() {
	RootCmd.AddCommand(statsCmd)
	statsCmd.PersistentFlags().StringVarP(&statsintreestr, "input", "i", "stdin", "Input tree")
	statsCmd.PersistentFlags().StringVarP(&statsoutfile, "output", "o", "stdout", "Output file")
}
