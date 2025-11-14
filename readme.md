# pubSEM: Publication-Ready Path Diagrams in R

pubSEM is a tool designed for the easy creation of publication-ready path diagrams. The tool is still in its (very)
early stages. See the roadmap below for future plans.

## Installing

> [!IMPORTANT]
> pubSEM requires a recent version of Go for compilation. Go can be installed here: https://go.dev/doc/install


To install the package in r, run the following line:

```r
remotes::install_github("dylanwglenn/pubSEM")
```

## Usage

> [!NOTE]
> pubSEM is designed only for plotting the output of lavaan models. If you want to interactively create lavaan models by
> visually connecting nodes, try the wonderful [lavaangui](https://github.com/karchjd/lavaangui/).

To use pubSEM, first create a diagram layout from a fitted lavaan model.

```r
pubSEM::sem_gui(fit = lavaan_fit, layout_name = "my-layout", standardized = FALSE) 
```

When you first see your model, it will be laid out essentially at random. Click and drag the nodes to arrange the
diagram as you would like. **Be sure to press "ctrl/cmd-S" to save your layout!**

To export your diagram to PDF, use the `export_diagram` function.
```r
pubSEM::export_diagram(layout_name = "my-layout", filename = "my-awsome-path-diagram")
```

<img width="720" height="480" alt="example path diagram" src="https://github.com/user-attachments/assets/00ad0560-8f78-45fe-b348-7a812e6b239d" />

## Motivation

Existing tools for the creation of path diagrams are either programmatic or interactive.

Programmatic solutions, such as [semPlot](https://cran.r-layout.org/web/packages/semPlot/index.html)
and [tidySEM](https://cjvanlissa.github.io/tidySEM/), have the benefit of producing diagrams that are easily
reproducible, meaning that once a layout of nodes is defined, the output automatically adapts to new data and can be
exported within the script (i.e. without having to export from an external editor). The downsides of these solutions
include the unintuitive programmatic creation of layouts and difficulty in maintaining legibility for complex diagrams.

Interactive solutions, including [lavaangui](https://github.com/karchjd/lavaangui/), have the benefit of being
interactive and intuitive. With these tools, users can easily layout and customize a path diagram. On the other hand,
existing solutions in this category lack the exactness and reproducibility of programmatic solutions. After the user
manually positions nodes and connections (often with some slight misalignment), the layout is not remembered. If you
want redo their model with new data or add another variable, all your previous work in laying out the model must
be redone.

## The pubSEM solution

pubSEM seeks to offer the intuitiveness of interactive plotting solutions along with the precision and reproducibility
of programmatic solutions. The tool features an interactive editor that can be used to precisely position nodes. The
positions of these nodes are then remembered — once a layout is defined and saved with the`semGUI()` function,
`export_path_diagram()` can be used to export a path diagram according to the layout previously defined without opening
the GUI.

## Why is it not on CRAN?

The pubSEM GUI is a standalone Go executable compiled from source. As Go is not within the CRAN standard toolchain and
CRAN understandably makes it difficult to distribute pre-built executables as part of an R package, this package
will live only on GitHub for the foreseeable future.

## Roadmap
- [ ] Proper toolbar
- [ ] Customization of covariance curves
- [ ] Multiple node selection
- [ ] Editing the visual names of variables
- [ ] Adjusting color and weight of elements
- [ ] Option to show confidence interval
- [ ] Adjusting the height of nodes
- [ ] Multi-line variable names
- [ ] Option to use serif or sans-serif fonts
- [ ] Option to grey out insignificant paths
- [ ] Changing dimensions for grid snapping
- [ ] support for multiple groups
- [ ] support for multi-level models
