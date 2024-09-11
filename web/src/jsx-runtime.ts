
export type Component = (props: Record<string, any>) => any;

export const jsx = {
  component(
    component: string | Component,
    props: Record<string, any> | null,
    ...children: any[]
  ) {
    if (!props) props = {};
    props.children = children.flat(Infinity);

    if (typeof component === "function") return component(props);

    if (typeof component === "string" && component === "null") {
      const fragment = document.createDocumentFragment();
      fragment.append(...props.children);
      return fragment;
    }

    const element = document.createElement(component);
    for (const [key, value] of Object.entries(props)) {
      if (key === "children") continue;
      else if (key === "className") element.setAttribute("class", value);
      else element.setAttribute(key, value);
    }

    element.append(...props.children);

    return element;
  },
};
