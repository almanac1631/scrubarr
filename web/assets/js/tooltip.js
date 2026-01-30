const nanosecondsToDays = (ns) => Math.floor(ns / 86_400_000_000_000);

document.addEventListener('DOMContentLoaded', () => {
    const tooltip = document.getElementById('global-tooltip');

    document.addEventListener('mouseover', (e) => {
        const trigger = e.target.closest('[data-tooltip]');

        if (!trigger) {
            tooltip.classList.add('hidden', 'opacity-0');
            return;
        }

        const tooltipKey = trigger.dataset.tooltip;

        if (tooltipKey === "torrent-info") {
            const torrentStatus = trigger.dataset.torrentStatus;
            tooltip.textContent = `Torrent link ${torrentStatus}`
            const trackerName = trigger.dataset.trackerName;
            if (trackerName !== "") {
                tooltip.textContent += ` (${trackerName})`
            }
        } else if (tooltipKey === "ratio-info") {
            const ratioStatus = trigger.dataset.ratioStatus;
            const ratio = formatFloatStr(trigger.dataset.ratio);
            const minRatio = formatFloatStr(trigger.dataset.minRatio);
            tooltip.textContent = `Ratio ${ratioStatus}`
            if (minRatio !== "-1" && ratio !== "-1") {
                tooltip.textContent += ` (${ratio}/${minRatio})`
            }
        } else if (tooltipKey === "age-info") {
            const ageStatus = trigger.dataset.ageStatus;
            const age = trigger.dataset.age;
            const minAge = trigger.dataset.minAge;
            tooltip.textContent = `Age ${ageStatus}`
            if (minAge !== "-1" && age !== "-1") {
                const minAgeDays = nanosecondsToDays(Number(minAge))
                const ageDays = nanosecondsToDays(Number(age));
                tooltip.textContent += ` (${ageDays}d/${minAgeDays}d)`
            }
        } else {
            return;
        }

        const rect = trigger.getBoundingClientRect();
        tooltip.classList.remove('hidden');
        const top = rect.top - tooltip.offsetHeight - 8; // 8px gap
        const left = rect.left + (rect.width / 2) - (tooltip.offsetWidth / 2);
        tooltip.style.top = `${top}px`;
        tooltip.style.left = `${left}px`;

        requestAnimationFrame(() => {
            tooltip.classList.remove('opacity-0');
        });
    });

    document.addEventListener('mouseout', (e) => {
        const trigger = e.target.closest('[data-tooltip]');
        if (trigger) {
            tooltip.classList.add('opacity-0');
        }
    });
});


function formatFloatStr(floatStr) {
    return new Intl.NumberFormat('en-US', {
        minimumFractionDigits: 0,
        maximumFractionDigits: 2
    }).format(floatStr)
}